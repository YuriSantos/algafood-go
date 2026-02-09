package model

import "github.com/shopspring/decimal"

// Produto represents a product from a restaurant
type Produto struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome          string          `gorm:"size:80;not null" json:"nome"`
	Descricao     string          `gorm:"size:255" json:"descricao"`
	Preco         decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"preco"`
	Ativo         bool            `gorm:"default:true" json:"ativo"`
	RestauranteID uint64          `gorm:"not null" json:"restauranteId"`
	Restaurante   Restaurante     `gorm:"foreignKey:RestauranteID" json:"-"`
}

func (Produto) TableName() string {
	return "produto"
}
