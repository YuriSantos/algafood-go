package event

import (
	"time"

	"github.com/shopspring/decimal"
)

// DomainEvent é a interface base para todos os eventos de domínio
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// BaseEvent contém campos comuns a todos os eventos
type BaseEvent struct {
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// PedidoConfirmadoEvent é emitido quando um pedido é confirmado
type PedidoConfirmadoEvent struct {
	BaseEvent
	PedidoCodigo    string          `json:"pedidoCodigo"`
	ClienteID       uint64          `json:"clienteId"`
	ClienteNome     string          `json:"clienteNome"`
	ClienteEmail    string          `json:"clienteEmail"`
	RestauranteID   uint64          `json:"restauranteId"`
	RestauranteNome string          `json:"restauranteNome"`
	ValorTotal      decimal.Decimal `json:"valorTotal"`
	DataConfirmacao time.Time       `json:"dataConfirmacao"`
}

func (e PedidoConfirmadoEvent) EventType() string {
	return "PedidoConfirmado"
}

func NewPedidoConfirmadoEvent(
	pedidoCodigo string,
	clienteID uint64,
	clienteNome string,
	clienteEmail string,
	restauranteID uint64,
	restauranteNome string,
	valorTotal decimal.Decimal,
	dataConfirmacao time.Time,
) PedidoConfirmadoEvent {
	return PedidoConfirmadoEvent{
		BaseEvent:       BaseEvent{Timestamp: time.Now()},
		PedidoCodigo:    pedidoCodigo,
		ClienteID:       clienteID,
		ClienteNome:     clienteNome,
		ClienteEmail:    clienteEmail,
		RestauranteID:   restauranteID,
		RestauranteNome: restauranteNome,
		ValorTotal:      valorTotal,
		DataConfirmacao: dataConfirmacao,
	}
}

// PedidoCanceladoEvent é emitido quando um pedido é cancelado
type PedidoCanceladoEvent struct {
	BaseEvent
	PedidoCodigo     string          `json:"pedidoCodigo"`
	ClienteID        uint64          `json:"clienteId"`
	ClienteNome      string          `json:"clienteNome"`
	ClienteEmail     string          `json:"clienteEmail"`
	RestauranteID    uint64          `json:"restauranteId"`
	RestauranteNome  string          `json:"restauranteNome"`
	ValorTotal       decimal.Decimal `json:"valorTotal"`
	DataCancelamento time.Time       `json:"dataCancelamento"`
}

func (e PedidoCanceladoEvent) EventType() string {
	return "PedidoCancelado"
}

func NewPedidoCanceladoEvent(
	pedidoCodigo string,
	clienteID uint64,
	clienteNome string,
	clienteEmail string,
	restauranteID uint64,
	restauranteNome string,
	valorTotal decimal.Decimal,
	dataCancelamento time.Time,
) PedidoCanceladoEvent {
	return PedidoCanceladoEvent{
		BaseEvent:        BaseEvent{Timestamp: time.Now()},
		PedidoCodigo:     pedidoCodigo,
		ClienteID:        clienteID,
		ClienteNome:      clienteNome,
		ClienteEmail:     clienteEmail,
		RestauranteID:    restauranteID,
		RestauranteNome:  restauranteNome,
		ValorTotal:       valorTotal,
		DataCancelamento: dataCancelamento,
	}
}

// PedidoEntregueEvent é emitido quando um pedido é entregue
type PedidoEntregueEvent struct {
	BaseEvent
	PedidoCodigo    string          `json:"pedidoCodigo"`
	ClienteID       uint64          `json:"clienteId"`
	ClienteNome     string          `json:"clienteNome"`
	ClienteEmail    string          `json:"clienteEmail"`
	RestauranteID   uint64          `json:"restauranteId"`
	RestauranteNome string          `json:"restauranteNome"`
	ValorTotal      decimal.Decimal `json:"valorTotal"`
	DataEntrega     time.Time       `json:"dataEntrega"`
}

func (e PedidoEntregueEvent) EventType() string {
	return "PedidoEntregue"
}

func NewPedidoEntregueEvent(
	pedidoCodigo string,
	clienteID uint64,
	clienteNome string,
	clienteEmail string,
	restauranteID uint64,
	restauranteNome string,
	valorTotal decimal.Decimal,
	dataEntrega time.Time,
) PedidoEntregueEvent {
	return PedidoEntregueEvent{
		BaseEvent:       BaseEvent{Timestamp: time.Now()},
		PedidoCodigo:    pedidoCodigo,
		ClienteID:       clienteID,
		ClienteNome:     clienteNome,
		ClienteEmail:    clienteEmail,
		RestauranteID:   restauranteID,
		RestauranteNome: restauranteNome,
		ValorTotal:      valorTotal,
		DataEntrega:     dataEntrega,
	}
}
