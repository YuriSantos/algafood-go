package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/yurisasc/algafood-go/internal/config"
)

// MessageHandler é a interface para processar mensagens do SQS
type MessageHandler interface {
	Handle(ctx context.Context, message *SQSMessage) error
}

// SQSMessage representa a estrutura de uma mensagem do EventBridge via SQS
type SQSMessage struct {
	Version    string          `json:"version"`
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	Account    string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Detail     json.RawMessage `json:"detail"`
}

// SQSListener escuta mensagens de uma fila SQS
type SQSListener struct {
	client            *sqs.Client
	queueURL          string
	handler           MessageHandler
	maxMessages       int32
	waitTimeSeconds   int32
	visibilityTimeout int32
	running           bool
	stopChan          chan struct{}
}

// NewSQSListener cria um novo listener SQS
func NewSQSListener(cfg *config.SQSConfig, awsCfg *config.AWSConfig, handler MessageHandler) (*SQSListener, error) {
	var opts []func(*awsconfig.LoadOptions) error

	opts = append(opts, awsconfig.WithRegion(cfg.Region))

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

	// Opções do cliente SQS
	var clientOpts []func(*sqs.Options)

	// Se tiver endpoint customizado (LocalStack), usa ele
	if awsCfg != nil && awsCfg.EndpointURL != "" {
		clientOpts = append(clientOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(awsCfg.EndpointURL)
		})
	}

	client := sqs.NewFromConfig(sdkCfg, clientOpts...)

	// Obter URL da fila
	queueURL := cfg.QueueURL
	if !isURL(queueURL) {
		// Se não for URL completa, tenta obter via GetQueueUrl
		log.Printf("Obtendo URL da fila para o nome: %s", queueURL)
		result, err := client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
			QueueName: aws.String(queueURL),
		})
		if err != nil {
			// Fallback: construir URL manualmente para LocalStack
			if awsCfg != nil && awsCfg.EndpointURL != "" {
				queueURL = fmt.Sprintf("%s/000000000000/%s", awsCfg.EndpointURL, cfg.QueueURL)
				log.Printf("Usando URL de fallback da fila: %s", queueURL)
			} else {
				return nil, fmt.Errorf("falha ao obter URL da fila: %w", err)
			}
		} else {
			queueURL = *result.QueueUrl
			log.Printf("URL da fila obtida da AWS: %s", queueURL)
		}
	}

	log.Printf("SQS Listener configurado - URL da Fila: %s, Região: %s", queueURL, cfg.Region)

	return &SQSListener{
		client:            client,
		queueURL:          queueURL,
		handler:           handler,
		maxMessages:       int32(cfg.MaxMessages),
		waitTimeSeconds:   int32(cfg.WaitTimeSeconds),
		visibilityTimeout: int32(cfg.VisibilityTimeout),
		stopChan:          make(chan struct{}),
	}, nil
}

// isURL verifica se a string é uma URL
func isURL(s string) bool {
	return len(s) > 7 && (s[:7] == "http://" || s[:8] == "https://")
}

// Start inicia o listener em uma goroutine
func (l *SQSListener) Start(ctx context.Context) {
	l.running = true
	log.Printf("Iniciando SQS listener para a fila: %s", l.queueURL)

	go func() {
		for l.running {
			select {
			case <-l.stopChan:
				log.Println("SQS listener parado")
				return
			case <-ctx.Done():
				log.Println("Contexto do SQS listener cancelado")
				return
			default:
				l.pollMessages(ctx)
			}
		}
	}()
}

// Stop para o listener
func (l *SQSListener) Stop() {
	l.running = false
	close(l.stopChan)
}

// pollMessages busca e processa mensagens da fila
func (l *SQSListener) pollMessages(ctx context.Context) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(l.queueURL),
		MaxNumberOfMessages: l.maxMessages,
		WaitTimeSeconds:     l.waitTimeSeconds,
		VisibilityTimeout:   l.visibilityTimeout,
		MessageAttributeNames: []string{
			"All",
		},
	}

	result, err := l.client.ReceiveMessage(ctx, input)
	if err != nil {
		log.Printf("Erro ao receber mensagens do SQS (fila: %s): %v", l.queueURL, err)
		time.Sleep(5 * time.Second) // Backoff em caso de erro
		return
	}

	if len(result.Messages) > 0 {
		log.Printf("Recebidas %d mensagens da fila SQS: %s", len(result.Messages), l.queueURL)
	}

	for _, msg := range result.Messages {
		log.Printf("Processando mensagem SQS ID: %s, Corpo: %s", *msg.MessageId, *msg.Body)

		if err := l.processMessage(ctx, msg); err != nil {
			log.Printf("Erro ao processar mensagem %s: %v", *msg.MessageId, err)
			// Não deleta a mensagem para que seja reprocessada
			continue
		}

		// Deleta a mensagem após processamento bem-sucedido
		if err := l.deleteMessage(ctx, msg.ReceiptHandle); err != nil {
			log.Printf("Erro ao deletar mensagem %s: %v", *msg.MessageId, err)
		} else {
			log.Printf("Mensagem processada e deletada com sucesso: %s", *msg.MessageId)
		}
	}
}

// processMessage processa uma mensagem individual
func (l *SQSListener) processMessage(ctx context.Context, msg types.Message) error {
	var sqsMessage SQSMessage
	if err := json.Unmarshal([]byte(*msg.Body), &sqsMessage); err != nil {
		return fmt.Errorf("falha ao deserializar mensagem: %w", err)
	}

	log.Printf("Processando mensagem: %s, Tipo de Detalhe: %s", *msg.MessageId, sqsMessage.DetailType)

	return l.handler.Handle(ctx, &sqsMessage)
}

// deleteMessage remove uma mensagem da fila
func (l *SQSListener) deleteMessage(ctx context.Context, receiptHandle *string) error {
	_, err := l.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(l.queueURL),
		ReceiptHandle: receiptHandle,
	})
	return err
}

// FakeSQSListener é um listener fake para desenvolvimento
type FakeSQSListener struct {
	handler MessageHandler
}

func NewFakeSQSListener(handler MessageHandler) *FakeSQSListener {
	return &FakeSQSListener{handler: handler}
}

func (l *FakeSQSListener) Start(ctx context.Context) {
	log.Println("[FAKE SQS] Listener iniciado (sem operação)")
}

func (l *FakeSQSListener) Stop() {
	log.Println("[FAKE SQS] Listener parado (sem operação)")
}

// SQSListenerInterface define a interface comum para listeners
type SQSListenerInterface interface {
	Start(ctx context.Context)
	Stop()
}

// NewSQSListenerFromConfig cria o listener baseado na configuração
func NewSQSListenerFromConfig(cfg *config.SQSConfig, awsCfg *config.AWSConfig, handler MessageHandler) (SQSListenerInterface, error) {
	if cfg.Type == "fake" {
		return NewFakeSQSListener(handler), nil
	}
	return NewSQSListener(cfg, awsCfg, handler)
}
