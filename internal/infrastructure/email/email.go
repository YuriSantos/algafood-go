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

// EmailMessage represents an email to be sent
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
}

// EmailService interface for sending emails
type EmailService interface {
	Send(message EmailMessage) error
}

// NewEmailService creates a new email service based on configuration
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

// FakeEmailService logs emails instead of sending
type FakeEmailService struct{}

func NewFakeEmailService() *FakeEmailService {
	return &FakeEmailService{}
}

func (s *FakeEmailService) Send(message EmailMessage) error {
	log.Printf("[FAKE EMAIL] To: %v, Subject: %s, Body: %s",
		message.To, message.Subject, message.Body)
	return nil
}

// SandboxEmailService sends all emails to a specific recipient
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
	// Override recipients with sandbox recipient
	sandboxMessage := EmailMessage{
		To:      []string{s.recipient},
		Subject: "[SANDBOX] " + message.Subject,
		Body:    fmt.Sprintf("Original recipients: %v\n\n%s", message.To, message.Body),
	}
	return s.smtpSvc.Send(sandboxMessage)
}

// SMTPEmailService sends emails via SendGrid
type SMTPEmailService struct {
	apiKey string
	from   string
}

func NewSMTPEmailService(cfg *config.EmailConfig) *SMTPEmailService {
	return &SMTPEmailService{
		apiKey: cfg.SMTP.Password, // SendGrid API key is stored as password
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
			return fmt.Errorf("failed to send email: %w", err)
		}

		if response.StatusCode >= 400 {
			return fmt.Errorf("email sending failed with status %d: %s",
				response.StatusCode, response.Body)
		}
	}

	return nil
}

// SESEmailService sends emails via AWS SES
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
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
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
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	log.Printf("Email sent via SES. MessageId: %s", *result.MessageId)
	return nil
}
