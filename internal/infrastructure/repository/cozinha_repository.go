package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type cozinhaRepositoryImpl struct {
	db *gorm.DB
}

// NewCozinhaRepository creates a new CozinhaRepository
func NewCozinhaRepository(db *gorm.DB) *cozinhaRepositoryImpl {
	return &cozinhaRepositoryImpl{db: db}
}

func (r *cozinhaRepositoryImpl) FindAll(page *pagination.Pageable) (*pagination.Page[model.Cozinha], error) {
	var cozinhas []model.Cozinha
	var total int64

	r.db.Model(&model.Cozinha{}).Count(&total)

	if err := r.db.Offset(page.Offset()).Limit(page.Size).Find(&cozinhas).Error; err != nil {
		return nil, err
	}

	return pagination.NewPage(cozinhas, total, page), nil
}

func (r *cozinhaRepositoryImpl) FindByID(id uint64) (*model.Cozinha, error) {
	var cozinha model.Cozinha
	if err := r.db.First(&cozinha, id).Error; err != nil {
		return nil, err
	}
	return &cozinha, nil
}

func (r *cozinhaRepositoryImpl) Save(cozinha *model.Cozinha) error {
	return r.db.Save(cozinha).Error
}

func (r *cozinhaRepositoryImpl) Delete(id uint64) error {
	return r.db.Delete(&model.Cozinha{}, id).Error
}
