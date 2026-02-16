package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/yurisasc/algafood-go/internal/infrastructure/email"
	"github.com/yurisasc/algafood-go/internal/infrastructure/sqs"
)

// PedidoEventDetail representa os dados do evento de pedido
type PedidoEventDetail struct {
	Timestamp        time.Time       `json:"timestamp"`
	PedidoCodigo     string          `json:"pedidoCodigo"`
	ClienteID        uint64          `json:"clienteId"`
	ClienteNome      string          `json:"clienteNome"`
	ClienteEmail     string          `json:"clienteEmail"`
	RestauranteID    uint64          `json:"restauranteId"`
	RestauranteNome  string          `json:"restauranteNome"`
	ValorTotal       decimal.Decimal `json:"valorTotal"`
	DataConfirmacao  *time.Time      `json:"dataConfirmacao,omitempty"`
	DataCancelamento *time.Time      `json:"dataCancelamento,omitempty"`
	DataEntrega      *time.Time      `json:"dataEntrega,omitempty"`
}

// NotificationHandler processa mensagens SQS e envia notificações por email
type NotificationHandler struct {
	emailService email.EmailService
}

// NewNotificationHandler cria um novo handler de notificações
func NewNotificationHandler(emailService email.EmailService) *NotificationHandler {
	return &NotificationHandler{
		emailService: emailService,
	}
}

// Handle processa uma mensagem SQS
func (h *NotificationHandler) Handle(ctx context.Context, message *sqs.SQSMessage) error {
	log.Printf("Handling notification for event: %s", message.DetailType)

	switch message.DetailType {
	case "PedidoConfirmado":
		return h.handlePedidoConfirmado(ctx, message.Detail)
	case "PedidoCancelado":
		return h.handlePedidoCancelado(ctx, message.Detail)
	case "PedidoEntregue":
		return h.handlePedidoEntregue(ctx, message.Detail)
	default:
		log.Printf("Unknown event type: %s", message.DetailType)
		return nil // Ignora eventos desconhecidos
	}
}

func (h *NotificationHandler) handlePedidoConfirmado(ctx context.Context, detail json.RawMessage) error {
	var evento PedidoEventDetail
	if err := json.Unmarshal(detail, &evento); err != nil {
		return fmt.Errorf("failed to unmarshal PedidoConfirmado: %w", err)
	}

	subject := fmt.Sprintf("Pedido confirmado - Código: %s", evento.PedidoCodigo)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Pedido Confirmado!</h1>
			<p>Olá, <strong>%s</strong>!</p>
			<p>Seu pedido no restaurante <strong>%s</strong> foi confirmado com sucesso.</p>
			<hr>
			<h3>Detalhes do Pedido:</h3>
			<ul>
				<li><strong>Código:</strong> %s</li>
				<li><strong>Valor Total:</strong> R$ %s</li>
				<li><strong>Data de Confirmação:</strong> %s</li>
			</ul>
			<hr>
			<p>Obrigado por escolher o AlgaFood!</p>
		</body>
		</html>
	`, evento.ClienteNome, evento.RestauranteNome, evento.PedidoCodigo,
		evento.ValorTotal.StringFixed(2), formatTime(evento.DataConfirmacao))

	return h.sendEmail(evento.ClienteEmail, subject, body)
}

func (h *NotificationHandler) handlePedidoCancelado(ctx context.Context, detail json.RawMessage) error {
	var evento PedidoEventDetail
	if err := json.Unmarshal(detail, &evento); err != nil {
		return fmt.Errorf("failed to unmarshal PedidoCancelado: %w", err)
	}

	subject := fmt.Sprintf("Pedido cancelado - Código: %s", evento.PedidoCodigo)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Pedido Cancelado</h1>
			<p>Olá, <strong>%s</strong>!</p>
			<p>Infelizmente, seu pedido no restaurante <strong>%s</strong> foi cancelado.</p>
			<hr>
			<h3>Detalhes do Pedido:</h3>
			<ul>
				<li><strong>Código:</strong> %s</li>
				<li><strong>Valor Total:</strong> R$ %s</li>
				<li><strong>Data de Cancelamento:</strong> %s</li>
			</ul>
			<hr>
			<p>Caso tenha dúvidas, entre em contato conosco.</p>
			<p>Equipe AlgaFood</p>
		</body>
		</html>
	`, evento.ClienteNome, evento.RestauranteNome, evento.PedidoCodigo,
		evento.ValorTotal.StringFixed(2), formatTime(evento.DataCancelamento))

	return h.sendEmail(evento.ClienteEmail, subject, body)
}

func (h *NotificationHandler) handlePedidoEntregue(ctx context.Context, detail json.RawMessage) error {
	var evento PedidoEventDetail
	if err := json.Unmarshal(detail, &evento); err != nil {
		return fmt.Errorf("failed to unmarshal PedidoEntregue: %w", err)
	}

	subject := fmt.Sprintf("Pedido entregue - Código: %s", evento.PedidoCodigo)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Pedido Entregue!</h1>
			<p>Olá, <strong>%s</strong>!</p>
			<p>Seu pedido do restaurante <strong>%s</strong> foi entregue com sucesso!</p>
			<hr>
			<h3>Detalhes do Pedido:</h3>
			<ul>
				<li><strong>Código:</strong> %s</li>
				<li><strong>Valor Total:</strong> R$ %s</li>
				<li><strong>Data de Entrega:</strong> %s</li>
			</ul>
			<hr>
			<p>Esperamos que aproveite sua refeição!</p>
			<p>Obrigado por escolher o AlgaFood!</p>
		</body>
		</html>
	`, evento.ClienteNome, evento.RestauranteNome, evento.PedidoCodigo,
		evento.ValorTotal.StringFixed(2), formatTime(evento.DataEntrega))

	return h.sendEmail(evento.ClienteEmail, subject, body)
}

func (h *NotificationHandler) sendEmail(to, subject, body string) error {
	message := email.EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	if err := h.emailService.Send(message); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func formatTime(t *time.Time) string {
	if t == nil {
		return "N/A"
	}
	return t.Format("02/01/2006 15:04:05")
}
