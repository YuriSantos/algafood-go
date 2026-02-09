package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type estadoRepositoryImpl struct {
	db *gorm.DB
}

// NewEstadoRepository creates a new EstadoRepository
func NewEstadoRepository(db *gorm.DB) *estadoRepositoryImpl {
	return &estadoRepositoryImpl{db: db}
}

func (r *estadoRepositoryImpl) FindAll() ([]model.Estado, error) {
	var estados []model.Estado
	if err := r.db.Find(&estados).Error; err != nil {
		return nil, err
	}
	return estados, nil
}

func (r *estadoRepositoryImpl) FindByID(id uint64) (*model.Estado, error) {
	var estado model.Estado
	if err := r.db.First(&estado, id).Error; err != nil {
		return nil, err
	}
	return &estado, nil
}

func (r *estadoRepositoryImpl) Save(estado *model.Estado) error {
	return r.db.Save(estado).Error
}

func (r *estadoRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&model.Estado{}, id).Error
}
