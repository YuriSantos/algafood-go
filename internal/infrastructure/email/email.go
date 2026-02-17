package email

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/yurisasc/algafood-go/internal/config"
)

// EmailMessage representa um email a ser enviado
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}

// EmailService interface para envio de emails
type EmailService interface {
	Send(message EmailMessage) error
}

// NewEmailService cria um novo serviço de email baseado na configuração
func NewEmailService(cfg *config.EmailConfig, awsCfg *config.AWSConfig) (EmailService, error) {
	switch cfg.Type {
	case "ses":
		return NewSESEmailService(cfg, awsCfg)
	case "smtp":
		return NewSMTPEmailService(cfg), nil
	case "sandbox":
		return NewSandboxEmailService(cfg), nil
	default:
		return NewFakeEmailService(), nil
	}
}

// FakeEmailService registra emails em log ao invés de enviar
type FakeEmailService struct{}

func NewFakeEmailService() *FakeEmailService {
	return &FakeEmailService{}
}

func (s *FakeEmailService) Send(message EmailMessage) error {
	log.Printf("[FAKE EMAIL] Para: %v, Assunto: %s, Corpo: %s",
		message.To, message.Subject, message.Body)
	return nil
}

// SandboxEmailService envia todos os emails para um destinatário específico
type SandboxEmailService struct {
	recipient string
	from      string
	smtpSvc   *SMTPEmailService
}

func NewSandboxEmailService(cfg *config.EmailConfig) *SandboxEmailService {
	return &SandboxEmailService{
		recipient: cfg.Sandbox.Recipient,
		from:      cfg.From,
		smtpSvc:   NewSMTPEmailService(cfg),
	}
}

func (s *SandboxEmailService) Send(message EmailMessage) error {
	// Sobrescreve destinatários com destinatário sandbox
	sandboxMessage := EmailMessage{
		To:      []string{s.recipient},
		Subject: "[SANDBOX] " + message.Subject,
		Body:    fmt.Sprintf("Destinatários originais: %v\n\n%s", message.To, message.Body),
	}
	return s.smtpSvc.Send(sandboxMessage)
}

// SMTPEmailService envia emails via SendGrid
type SMTPEmailService struct {
	apiKey string
	from   string
}

func NewSMTPEmailService(cfg *config.EmailConfig) *SMTPEmailService {
	return &SMTPEmailService{
		apiKey: cfg.SMTP.Password, // Chave API do SendGrid armazenada como password
		from:   cfg.From,
	}
}

func (s *SMTPEmailService) Send(message EmailMessage) error {
	from := mail.NewEmail("AlgaFood", s.from)

	for _, recipient := range message.To {
		to := mail.NewEmail("", recipient)
		email := mail.NewSingleEmail(from, message.Subject, to, message.Body, message.Body)
		client := sendgrid.NewSendClient(s.apiKey)

		response, err := client.Send(email)
		if err != nil {
			return fmt.Errorf("falha ao enviar email: %w", err)
		}

		if response.StatusCode >= 400 {
			return fmt.Errorf("falha ao enviar email com status %d: %s",
				response.StatusCode, response.Body)
		}
	}

	return nil
}

// SESEmailService envia emails via AWS SES
type SESEmailService struct {
	client *ses.Client
	from   string
}

func NewSESEmailService(cfg *config.EmailConfig, awsCfg *config.AWSConfig) (*SESEmailService, error) {
	var opts []func(*awsconfig.LoadOptions) error

	opts = append(opts, awsconfig.WithRegion(cfg.SES.Region))

	// Se tiver credenciais configuradas (LocalStack), usa elas
	if awsCfg != nil && awsCfg.Credentials.AccessKey != "" {
		opts = append(opts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				awsCfg.Credentials.AccessKey,
				awsCfg.Credentials.SecretKey,
				"",
			),
		))
	}

	sdkCfg, err := awsconfig.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar configuração AWS: %w", err)
	}

	// Opções do cliente SES
	var clientOpts []func(*ses.Options)

	// Se tiver endpoint customizado (LocalStack), usa ele
	if awsCfg != nil && awsCfg.EndpointURL != "" {
		clientOpts = append(clientOpts, func(o *ses.Options) {
			o.BaseEndpoint = aws.String(awsCfg.EndpointURL)
		})
	}

	client := ses.NewFromConfig(sdkCfg, clientOpts...)

	return &SESEmailService{
		client: client,
		from:   cfg.From,
	}, nil
}

func (s *SESEmailService) Send(message EmailMessage) error {
	ctx := context.Background()

	toAddresses := make([]string, len(message.To))
	copy(toAddresses, message.To)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: toAddresses,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message.Body),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message.Body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(message.Subject),
			},
		},
		Source: aws.String(s.from),
	}

	result, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("falha ao enviar email via SES: %w", err)
	}

	log.Printf("Email enviado via SES. MessageId: %s", *result.MessageId)
	return nil
}
