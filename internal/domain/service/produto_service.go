package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type ProdutoService struct {
	repo           repository.ProdutoRepository
	restauranteSvc *RestauranteService
}

func NewProdutoService(repo repository.ProdutoRepository, restauranteSvc *RestauranteService) *ProdutoService {
	return &ProdutoService{
		repo:           repo,
		restauranteSvc: restauranteSvc,
	}
}

func (s *ProdutoService) FindAllByRestaurante(restauranteID uint64, incluirInativos bool) ([]model.Produto, error) {
	// Validate restaurante exists
	if _, err := s.restauranteSvc.FindByID(restauranteID); err != nil {
		return nil, err
	}
	return s.repo.FindAllByRestaurante(restauranteID, incluirInativos)
}

func (s *ProdutoService) FindByID(restauranteID, produtoID uint64) (*model.Produto, error) {
	// Validate restaurante exists
	if _, err := s.restauranteSvc.FindByID(restauranteID); err != nil {
		return nil, err
	}

	produto, err := s.repo.FindByID(restauranteID, produtoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewProdutoNaoEncontradoException(restauranteID, produtoID)
		}
		return nil, err
	}
	return produto, nil
}

func (s *ProdutoService) Save(restauranteID uint64, produto *model.Produto) error {
	// Validate restaurante exists
	if _, err := s.restauranteSvc.FindByID(restauranteID); err != nil {
		return err
	}

	produto.RestauranteID = restauranteID
	return s.repo.Save(produto)
}
