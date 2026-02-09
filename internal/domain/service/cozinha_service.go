package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type CozinhaService struct {
	repo repository.CozinhaRepository
}

func NewCozinhaService(repo repository.CozinhaRepository) *CozinhaService {
	return &CozinhaService{repo: repo}
}

func (s *CozinhaService) FindAll(page *pagination.Pageable) (*pagination.Page[model.Cozinha], error) {
	return s.repo.FindAll(page)
}

func (s *CozinhaService) FindByID(id uint64) (*model.Cozinha, error) {
	cozinha, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewCozinhaNaoEncontradaException(id)
		}
		return nil, err
	}
	return cozinha, nil
}

func (s *CozinhaService) Save(cozinha *model.Cozinha) error {
	return s.repo.Save(cozinha)
}

func (s *CozinhaService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Cozinha nao pode ser removida, pois esta em uso")
	}
	return nil
}
