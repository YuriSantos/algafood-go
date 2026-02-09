package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"gorm.io/gorm"
)

type produtoRepositoryImpl struct {
	db *gorm.DB
}

// NewProdutoRepository creates a new ProdutoRepository
func NewProdutoRepository(db *gorm.DB) *produtoRepositoryImpl {
	return &produtoRepositoryImpl{db: db}
}

func (r *produtoRepositoryImpl) FindAllByRestaurante(restauranteID uint64, incluirInativos bool) ([]model.Produto, error) {
	var produtos []model.Produto
	query := r.db.Where("restaurante_id = ?", restauranteID)

	if !incluirInativos {
		query = query.Where("ativo = ?", true)
	}

	if err := query.Find(&produtos).Error; err != nil {
		return nil, err
	}
	return produtos, nil
}

func (r *produtoRepositoryImpl) FindByID(restauranteID, produtoID uint64) (*model.Produto, error) {
	var produto model.Produto
	if err := r.db.Where("restaurante_id = ? AND id = ?", restauranteID, produtoID).First(&produto).Error; err != nil {
		return nil, err
	}
	return &produto, nil
}

func (r *produtoRepositoryImpl) Save(produto *model.Produto) error {
	return r.db.Save(produto).Error
}
