package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type fotoProdutoRepositoryImpl struct {
	db *gorm.DB
}

// NewFotoProdutoRepository creates a new FotoProdutoRepository
func NewFotoProdutoRepository(db *gorm.DB) *fotoProdutoRepositoryImpl {
	return &fotoProdutoRepositoryImpl{db: db}
}

func (r *fotoProdutoRepositoryImpl) FindByProdutoID(produtoID uint64) (*model.FotoProduto, error) {
	var foto model.FotoProduto
	if err := r.db.Where("produto_id = ?", produtoID).First(&foto).Error; err != nil {
		return nil, err
	}
	return &foto, nil
}

func (r *fotoProdutoRepositoryImpl) Save(foto *model.FotoProduto) error {
	// Use upsert - update if exists, insert if not
	return r.db.Save(foto).Error
}

func (r *fotoProdutoRepositoryImpl) Delete(produtoID uint64) error {
	return r.db.Where("produto_id = ?", produtoID).Delete(&model.FotoProduto{}).Error
}
