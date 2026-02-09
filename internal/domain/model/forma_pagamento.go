package model

import "time"

// FormaPagamento represents a payment method
type FormaPagamento struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Descricao       string    `gorm:"size:60;not null" json:"descricao"`
	DataAtualizacao time.Time `gorm:"autoUpdateTime" json:"dataAtualizacao"`
}

func (FormaPagamento) TableName() string {
	return "forma_pagamento"
}
