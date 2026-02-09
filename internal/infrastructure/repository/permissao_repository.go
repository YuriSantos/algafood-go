package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type permissaoRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissaoRepository creates a new PermissaoRepository
func NewPermissaoRepository(db *gorm.DB) *permissaoRepositoryImpl {
	return &permissaoRepositoryImpl{db: db}
}

func (r *permissaoRepositoryImpl) FindAll() ([]model.Permissao, error) {
	var permissoes []model.Permissao
	if err := r.db.Find(&permissoes).Error; err != nil {
		return nil, err
	}
	return permissoes, nil
}

func (r *permissaoRepositoryImpl) FindByID(id uint64) (*model.Permissao, error) {
	var permissao model.Permissao
	if err := r.db.First(&permissao, id).Error; err != nil {
		return nil, err
	}
	return &permissao, nil
}
