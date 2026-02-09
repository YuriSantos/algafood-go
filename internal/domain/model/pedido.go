package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Pedido represents an order
type Pedido struct {
	ID               uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo           string          `gorm:"size:36;uniqueIndex;not null" json:"codigo"`
	Subtotal         decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	TaxaFrete        decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"taxaFrete"`
	ValorTotal       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"valorTotal"`
	Status           StatusPedido    `gorm:"type:varchar(15);not null;default:'CRIADO'" json:"status"`
	DataCriacao      time.Time       `gorm:"autoCreateTime" json:"dataCriacao"`
	DataConfirmacao  *time.Time      `json:"dataConfirmacao,omitempty"`
	DataCancelamento *time.Time      `json:"dataCancelamento,omitempty"`
	DataEntrega      *time.Time      `json:"dataEntrega,omitempty"`

	// Foreign keys
	RestauranteID    uint64         `gorm:"not null" json:"restauranteId"`
	Restaurante      Restaurante    `gorm:"foreignKey:RestauranteID" json:"restaurante,omitempty"`
	ClienteID        uint64         `gorm:"column:usuario_cliente_id;not null" json:"clienteId"`
	Cliente          Usuario        `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	FormaPagamentoID uint64         `gorm:"not null" json:"formaPagamentoId"`
	FormaPagamento   FormaPagamento `gorm:"foreignKey:FormaPagamentoID" json:"formaPagamento,omitempty"`

	// Embedded address
	EnderecoEntrega EnderecoEntrega `gorm:"embedded" json:"enderecoEntrega,omitempty"`

	// Items
	Itens []ItemPedido `gorm:"foreignKey:PedidoID" json:"itens,omitempty"`
}

func (Pedido) TableName() string {
	return "pedido"
}

// EnderecoEntrega represents the delivery address (embedded in Pedido)
type EnderecoEntrega struct {
	CEP         string `gorm:"column:endereco_cep;size:9" json:"cep"`
	Logradouro  string `gorm:"column:endereco_logradouro;size:100" json:"logradouro"`
	Numero      string `gorm:"column:endereco_numero;size:20" json:"numero"`
	Complemento string `gorm:"column:endereco_complemento;size:60" json:"complemento"`
	Bairro      string `gorm:"column:endereco_bairro;size:60" json:"bairro"`
	CidadeID    uint64 `gorm:"column:endereco_cidade_id" json:"cidadeId"`
	Cidade      Cidade `gorm:"foreignKey:CidadeID" json:"cidade,omitempty"`
}

// BeforeCreate generates a UUID for the order code
func (p *Pedido) BeforeCreate() {
	if p.Codigo == "" {
		p.Codigo = uuid.New().String()
	}
	p.Status = StatusPedidoCriado
}

// CalcularValorTotal calculates the total order value
func (p *Pedido) CalcularValorTotal() {
	p.Subtotal = decimal.Zero
	for _, item := range p.Itens {
		item.CalcularPrecoTotal()
		p.Subtotal = p.Subtotal.Add(item.PrecoTotal)
	}
	p.ValorTotal = p.Subtotal.Add(p.TaxaFrete)
}

// DefinirFrete sets the freight rate
func (p *Pedido) DefinirFrete(taxaFrete decimal.Decimal) {
	p.TaxaFrete = taxaFrete
}

// AtribuirPedidoAosItens associates this order to all items
func (p *Pedido) AtribuirPedidoAosItens() {
	for i := range p.Itens {
		p.Itens[i].PedidoID = p.ID
	}
}

// Confirmar confirms the order
func (p *Pedido) Confirmar() error {
	if !p.Status.CanTransitionTo(StatusPedidoConfirmado) {
		return newStatusChangeError(p.Status, StatusPedidoConfirmado)
	}
	p.Status = StatusPedidoConfirmado
	now := time.Now()
	p.DataConfirmacao = &now
	return nil
}

// Entregar marks the order as delivered
func (p *Pedido) Entregar() error {
	if !p.Status.CanTransitionTo(StatusPedidoEntregue) {
		return newStatusChangeError(p.Status, StatusPedidoEntregue)
	}
	p.Status = StatusPedidoEntregue
	now := time.Now()
	p.DataEntrega = &now
	return nil
}

// Cancelar cancels the order
func (p *Pedido) Cancelar() error {
	if !p.Status.CanTransitionTo(StatusPedidoCancelado) {
		return newStatusChangeError(p.Status, StatusPedidoCancelado)
	}
	p.Status = StatusPedidoCancelado
	now := time.Now()
	p.DataCancelamento = &now
	return nil
}

// StatusChangeError represents an invalid status transition
type StatusChangeError struct {
	CurrentStatus StatusPedido
	TargetStatus  StatusPedido
}

func (e *StatusChangeError) Error() string {
	return "Status do pedido " + string(e.CurrentStatus) + " nao pode ser alterado para " + string(e.TargetStatus)
}

func newStatusChangeError(current, target StatusPedido) *StatusChangeError {
	return &StatusChangeError{
		CurrentStatus: current,
		TargetStatus:  target,
	}
}

// PodeSerConfirmado checks if order can be confirmed
func (p *Pedido) PodeSerConfirmado() bool {
	return p.Status.CanTransitionTo(StatusPedidoConfirmado)
}

// PodeSerEntregue checks if order can be delivered
func (p *Pedido) PodeSerEntregue() bool {
	return p.Status.CanTransitionTo(StatusPedidoEntregue)
}

// PodeSerCancelado checks if order can be cancelled
func (p *Pedido) PodeSerCancelado() bool {
	return p.Status.CanTransitionTo(StatusPedidoCancelado)
}
