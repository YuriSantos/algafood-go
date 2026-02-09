package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type restauranteRepositoryImpl struct {
	db *gorm.DB
}

// NewRestauranteRepository creates a new RestauranteRepository
func NewRestauranteRepository(db *gorm.DB) *restauranteRepositoryImpl {
	return &restauranteRepositoryImpl{db: db}
}

func (r *restauranteRepositoryImpl) FindAll() ([]model.Restaurante, error) {
	var restaurantes []model.Restaurante
	if err := r.db.Preload("Cozinha").Find(&restaurantes).Error; err != nil {
		return nil, err
	}
	return restaurantes, nil
}

func (r *restauranteRepositoryImpl) FindByID(id uint64) (*model.Restaurante, error) {
	var restaurante model.Restaurante
	if err := r.db.Preload("Cozinha").
		Preload("FormasPagamento").
		Preload("Responsaveis").
		Preload("Endereco.Cidade").
		Preload("Endereco.Cidade.Estado").
		First(&restaurante, id).Error; err != nil {
		return nil, err
	}
	return &restaurante, nil
}

func (r *restauranteRepositoryImpl) Save(restaurante *model.Restaurante) error {
	return r.db.Save(restaurante).Error
}

func (r *restauranteRepositoryImpl) AddFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	return r.db.Exec("INSERT INTO restaurante_forma_pagamento (restaurante_id, forma_pagamento_id) VALUES (?, ?)", restauranteID, formaPagamentoID).Error
}

func (r *restauranteRepositoryImpl) RemoveFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	return r.db.Exec("DELETE FROM restaurante_forma_pagamento WHERE restaurante_id = ? AND forma_pagamento_id = ?", restauranteID, formaPagamentoID).Error
}

func (r *restauranteRepositoryImpl) AddResponsavel(restauranteID, usuarioID uint64) error {
	return r.db.Exec("INSERT INTO restaurante_usuario_responsavel (restaurante_id, usuario_id) VALUES (?, ?)", restauranteID, usuarioID).Error
}

func (r *restauranteRepositoryImpl) RemoveResponsavel(restauranteID, usuarioID uint64) error {
	return r.db.Exec("DELETE FROM restaurante_usuario_responsavel WHERE restaurante_id = ? AND usuario_id = ?", restauranteID, usuarioID).Error
}

func (r *restauranteRepositoryImpl) ExistsResponsavel(restauranteID, usuarioID uint64) (bool, error) {
	var count int64
	if err := r.db.Table("restaurante_usuario_responsavel").
		Where("restaurante_id = ? AND usuario_id = ?", restauranteID, usuarioID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
