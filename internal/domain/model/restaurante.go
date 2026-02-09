package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// Restaurante represents a restaurant
type Restaurante struct {
	ID              uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome            string           `gorm:"size:80;not null" json:"nome"`
	TaxaFrete       decimal.Decimal  `gorm:"type:decimal(10,2);not null" json:"taxaFrete"`
	CozinhaID       uint64           `gorm:"not null" json:"cozinhaId"`
	Cozinha         Cozinha          `gorm:"foreignKey:CozinhaID" json:"cozinha,omitempty"`
	Endereco        Endereco         `gorm:"embedded" json:"endereco,omitempty"`
	Ativo           bool             `gorm:"default:true" json:"ativo"`
	Aberto          bool             `gorm:"default:false" json:"aberto"`
	DataCadastro    time.Time        `gorm:"autoCreateTime" json:"dataCadastro"`
	DataAtualizacao time.Time        `gorm:"autoUpdateTime" json:"dataAtualizacao"`
	FormasPagamento []FormaPagamento `gorm:"many2many:restaurante_forma_pagamento;" json:"formasPagamento,omitempty"`
	Responsaveis    []Usuario        `gorm:"many2many:restaurante_usuario_responsavel;" json:"responsaveis,omitempty"`
	Produtos        []Produto        `gorm:"foreignKey:RestauranteID" json:"produtos,omitempty"`
}

func (Restaurante) TableName() string {
	return "restaurante"
}

// Ativar activates the restaurant
func (r *Restaurante) Ativar() {
	r.Ativo = true
}

// Inativar deactivates the restaurant
func (r *Restaurante) Inativar() {
	r.Ativo = false
}

// Abrir opens the restaurant for orders
func (r *Restaurante) Abrir() {
	r.Aberto = true
}

// Fechar closes the restaurant
func (r *Restaurante) Fechar() {
	r.Aberto = false
}

// AdicionarFormaPagamento adds a payment method
func (r *Restaurante) AdicionarFormaPagamento(formaPagamento FormaPagamento) {
	r.FormasPagamento = append(r.FormasPagamento, formaPagamento)
}

// RemoverFormaPagamento removes a payment method
func (r *Restaurante) RemoverFormaPagamento(formaPagamento FormaPagamento) {
	for i, fp := range r.FormasPagamento {
		if fp.ID == formaPagamento.ID {
			r.FormasPagamento = append(r.FormasPagamento[:i], r.FormasPagamento[i+1:]...)
			return
		}
	}
}

// AdicionarResponsavel adds a responsible user
func (r *Restaurante) AdicionarResponsavel(usuario Usuario) {
	r.Responsaveis = append(r.Responsaveis, usuario)
}

// RemoverResponsavel removes a responsible user
func (r *Restaurante) RemoverResponsavel(usuario Usuario) {
	for i, u := range r.Responsaveis {
		if u.ID == usuario.ID {
			r.Responsaveis = append(r.Responsaveis[:i], r.Responsaveis[i+1:]...)
			return
		}
	}
}

// AceitaFormaPagamento checks if restaurant accepts the given payment method
func (r *Restaurante) AceitaFormaPagamento(formaPagamento FormaPagamento) bool {
	for _, fp := range r.FormasPagamento {
		if fp.ID == formaPagamento.ID {
			return true
		}
	}
	return false
}

// NaoAceitaFormaPagamento checks if restaurant doesn't accept the given payment method
func (r *Restaurante) NaoAceitaFormaPagamento(formaPagamento FormaPagamento) bool {
	return !r.AceitaFormaPagamento(formaPagamento)
}
