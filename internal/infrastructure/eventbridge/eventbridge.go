package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yurisasc/algafood-go/internal/config"
	"github.com/yurisasc/algafood-go/internal/domain/event"
)

// EventPublisher é a interface para publicação de eventos
type EventPublisher interface {
	Publish(ctx context.Context, event event.DomainEvent) error
}

// SQSEventMessage representa a estrutura de uma mensagem do EventBridge enviada via SQS
type SQSEventMessage struct {
	Version    string          `json:"version"`
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	Account    string          `json:"account"`
	Time       time.Time       `json:"time"`
	Region     string          `json:"region"`
	Detail     json.RawMessage `json:"detail"`
}

// EventBridgePublisher publica eventos no AWS EventBridge
type EventBridgePublisher struct {
	client       *eventbridge.Client
	sqsClient    *sqs.Client
	sqsQueueURL  string
	eventBusName string
	source       string
	region       string
	directSQS    bool // Quando true, envia diretamente para SQS (útil para LocalStack)
}

// NewEventBridgePublisher cria um novo publicador de eventos para o EventBridge
func NewEventBridgePublisher(cfg *config.EventBridgeConfig, sqsCfg *config.SQSConfig, awsCfg *config.AWSConfig) (*EventBridgePublisher, error) {
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
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Opções do cliente EventBridge
	var ebClientOpts []func(*eventbridge.Options)
	var sqsClientOpts []func(*sqs.Options)

	// Se tiver endpoint customizado (LocalStack), usa ele
	directSQS := false
	if awsCfg != nil && awsCfg.EndpointURL != "" {
		ebClientOpts = append(ebClientOpts, func(o *eventbridge.Options) {
			o.BaseEndpoint = aws.String(awsCfg.EndpointURL)
		})
		sqsClientOpts = append(sqsClientOpts, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(awsCfg.EndpointURL)
		})
		// Habilita envio direto para SQS no LocalStack
		directSQS = true
	}

	ebClient := eventbridge.NewFromConfig(sdkCfg, ebClientOpts...)
	sqsClient := sqs.NewFromConfig(sdkCfg, sqsClientOpts...)

	// Obter URL da fila SQS
	var sqsQueueURL string
	if sqsCfg != nil && sqsCfg.QueueURL != "" {
		if len(sqsCfg.QueueURL) > 7 && (sqsCfg.QueueURL[:7] == "http://" || sqsCfg.QueueURL[:8] == "https://") {
			sqsQueueURL = sqsCfg.QueueURL
		} else {
			// Tenta obter a URL da fila
			result, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
				QueueName: aws.String(sqsCfg.QueueURL),
			})
			if err != nil {
				log.Printf("Warning: Could not get SQS queue URL: %v", err)
				if awsCfg != nil && awsCfg.EndpointURL != "" {
					sqsQueueURL = fmt.Sprintf("%s/000000000000/%s", awsCfg.EndpointURL, sqsCfg.QueueURL)
				}
			} else {
				sqsQueueURL = *result.QueueUrl
			}
		}
		log.Printf("EventBridge publisher configured with SQS queue: %s (direct: %v)", sqsQueueURL, directSQS)
	}

	return &EventBridgePublisher{
		client:       ebClient,
		sqsClient:    sqsClient,
		sqsQueueURL:  sqsQueueURL,
		eventBusName: cfg.EventBusName,
		source:       cfg.Source,
		region:       cfg.Region,
		directSQS:    directSQS,
	}, nil
}

// Publish publica um evento de domínio no EventBridge
func (p *EventBridgePublisher) Publish(ctx context.Context, domainEvent event.DomainEvent) error {
	detail, err := json.Marshal(domainEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	log.Printf("Publishing event to EventBridge - Bus: %s, Source: %s, Type: %s",
		p.eventBusName, p.source, domainEvent.EventType())
	log.Printf("Event detail: %s", string(detail))

	input := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String(p.eventBusName),
				Source:       aws.String(p.source),
				DetailType:   aws.String(domainEvent.EventType()),
				Detail:       aws.String(string(detail)),
				Time:         aws.Time(domainEvent.OccurredAt()),
			},
		},
	}

	result, err := p.client.PutEvents(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish event to EventBridge: %w", err)
	}

	if result.FailedEntryCount > 0 {
		for _, entry := range result.Entries {
			if entry.ErrorCode != nil {
				return fmt.Errorf("failed to publish event: %s - %s", *entry.ErrorCode, *entry.ErrorMessage)
			}
		}
	}

	eventID := "unknown"
	if len(result.Entries) > 0 && result.Entries[0].EventId != nil {
		eventID = *result.Entries[0].EventId
	}

	log.Printf("Event published successfully to EventBridge: %s (EventId: %s)",
		domainEvent.EventType(), eventID)

	// Se directSQS estiver habilitado (LocalStack), envia também diretamente para SQS
	if p.directSQS && p.sqsQueueURL != "" {
		if err := p.publishToSQS(ctx, domainEvent, detail, eventID); err != nil {
			log.Printf("Warning: Failed to publish event directly to SQS: %v", err)
			// Não retorna erro pois o evento já foi publicado no EventBridge
		}
	}

	return nil
}

// publishToSQS envia o evento diretamente para a fila SQS
func (p *EventBridgePublisher) publishToSQS(ctx context.Context, domainEvent event.DomainEvent, detail []byte, eventID string) error {
	// Criar mensagem no formato do EventBridge
	sqsMessage := SQSEventMessage{
		Version:    "0",
		ID:         eventID,
		DetailType: domainEvent.EventType(),
		Source:     p.source,
		Account:    "000000000000",
		Time:       domainEvent.OccurredAt(),
		Region:     p.region,
		Detail:     detail,
	}

	messageBody, err := json.Marshal(sqsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal SQS message: %w", err)
	}

	_, err = p.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.sqsQueueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %w", err)
	}

	log.Printf("Event also published directly to SQS: %s (Queue: %s)", domainEvent.EventType(), p.sqsQueueURL)
	return nil
}

// FakeEventPublisher é um publicador de eventos para desenvolvimento/testes
type FakeEventPublisher struct{}

func NewFakeEventPublisher() *FakeEventPublisher {
	return &FakeEventPublisher{}
}

func (p *FakeEventPublisher) Publish(ctx context.Context, domainEvent event.DomainEvent) error {
	detail, _ := json.MarshalIndent(domainEvent, "", "  ")
	log.Printf("[FAKE EVENT] Type: %s, Time: %s, Detail:\n%s",
		domainEvent.EventType(),
		domainEvent.OccurredAt().Format(time.RFC3339),
		string(detail))
	return nil
}

// NewEventPublisher cria o publicador de eventos baseado na configuração
func NewEventPublisher(cfg *config.EventBridgeConfig, sqsCfg *config.SQSConfig, awsCfg *config.AWSConfig) (EventPublisher, error) {
	if cfg.Type == "fake" {
		return NewFakeEventPublisher(), nil
	}
	return NewEventBridgePublisher(cfg, sqsCfg, awsCfg)
}
