package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type cidadeRepositoryImpl struct {
	db *gorm.DB
}

// NewCidadeRepository creates a new CidadeRepository
func NewCidadeRepository(db *gorm.DB) *cidadeRepositoryImpl {
	return &cidadeRepositoryImpl{db: db}
}

func (r *cidadeRepositoryImpl) FindAll() ([]model.Cidade, error) {
	var cidades []model.Cidade
	if err := r.db.Preload("Estado").Find(&cidades).Error; err != nil {
		return nil, err
	}
	return cidades, nil
}

func (r *cidadeRepositoryImpl) FindByID(id uint64) (*model.Cidade, error) {
	var cidade model.Cidade
	if err := r.db.Preload("Estado").First(&cidade, id).Error; err != nil {
		return nil, err
	}
	return &cidade, nil
}

func (r *cidadeRepositoryImpl) Save(cidade *model.Cidade) error {
	return r.db.Save(cidade).Error
}

func (r *cidadeRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&model.Cidade{}, id).Error
}
