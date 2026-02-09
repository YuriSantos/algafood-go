package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type formaPagamentoRepositoryImpl struct {
	db *gorm.DB
}

// NewFormaPagamentoRepository creates a new FormaPagamentoRepository
func NewFormaPagamentoRepository(db *gorm.DB) *formaPagamentoRepositoryImpl {
	return &formaPagamentoRepositoryImpl{db: db}
}

func (r *formaPagamentoRepositoryImpl) FindAll() ([]model.FormaPagamento, error) {
	var formasPagamento []model.FormaPagamento
	if err := r.db.Find(&formasPagamento).Error; err != nil {
		return nil, err
	}
	return formasPagamento, nil
}

func (r *formaPagamentoRepositoryImpl) FindByID(id uint64) (*model.FormaPagamento, error) {
	var formaPagamento model.FormaPagamento
	if err := r.db.First(&formaPagamento, id).Error; err != nil {
		return nil, err
	}
	return &formaPagamento, nil
}

func (r *formaPagamentoRepositoryImpl) Save(formaPagamento *model.FormaPagamento) error {
	return r.db.Save(formaPagamento).Error
}

func (r *formaPagamentoRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&model.FormaPagamento{}, id).Error
}
