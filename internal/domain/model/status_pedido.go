package model

// StatusPedido represents the order status
type StatusPedido string

const (
	StatusPedidoCriado     StatusPedido = "CRIADO"
	StatusPedidoConfirmado StatusPedido = "CONFIRMADO"
	StatusPedidoEntregue   StatusPedido = "ENTREGUE"
	StatusPedidoCancelado  StatusPedido = "CANCELADO"
)

// CanTransitionTo checks if the current status can transition to the target status
func (s StatusPedido) CanTransitionTo(target StatusPedido) bool {
	transitions := map[StatusPedido][]StatusPedido{
		StatusPedidoCriado:     {StatusPedidoConfirmado, StatusPedidoCancelado},
		StatusPedidoConfirmado: {StatusPedidoEntregue, StatusPedidoCancelado},
	}

	allowedTargets, exists := transitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedTargets {
		if allowed == target {
			return true
		}
	}
	return false
}

// GetDescription returns the Portuguese description of the status
func (s StatusPedido) GetDescription() string {
	descriptions := map[StatusPedido]string{
		StatusPedidoCriado:     "Criado",
		StatusPedidoConfirmado: "Confirmado",
		StatusPedidoEntregue:   "Entregue",
		StatusPedidoCancelado:  "Cancelado",
	}
	return descriptions[s]
}
