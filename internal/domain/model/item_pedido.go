package model

import "github.com/shopspring/decimal"

// ItemPedido represents an order item
type ItemPedido struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PedidoID      uint64          `gorm:"not null" json:"pedidoId"`
	ProdutoID     uint64          `gorm:"not null" json:"produtoId"`
	Produto       Produto         `gorm:"foreignKey:ProdutoID" json:"produto,omitempty"`
	Quantidade    int             `gorm:"not null" json:"quantidade"`
	PrecoUnitario decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"precoUnitario"`
	PrecoTotal    decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"precoTotal"`
	Observacao    string          `gorm:"size:255" json:"observacao"`
}

func (ItemPedido) TableName() string {
	return "item_pedido"
}

// CalcularPrecoTotal calculates the total price of this item
func (i *ItemPedido) CalcularPrecoTotal() {
	i.PrecoTotal = i.PrecoUnitario.Mul(decimal.NewFromInt(int64(i.Quantidade)))
}
