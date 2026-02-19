package service

import (
	"errors"
	"log"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type CozinhaService struct {
	repo     repository.CozinhaRepository
	cacheSvc *BusinessCacheService
}

func NewCozinhaService(repo repository.CozinhaRepository, cacheSvc *BusinessCacheService) *CozinhaService {
	return &CozinhaService{
		repo:     repo,
		cacheSvc: cacheSvc,
	}
}

func (s *CozinhaService) FindAll(page *pagination.Pageable) (*pagination.Page[model.Cozinha], error) {
	return s.repo.FindAll(page)
}

func (s *CozinhaService) FindByID(id uint64) (*model.Cozinha, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetCozinha(id); err == nil && cached != nil {
			return cached, nil
		}
	}

	cozinha, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewCozinhaNaoEncontradaException(id)
		}
		return nil, err
	}

	// Armazena no cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetCozinha(cozinha); err != nil {
			log.Printf("Aviso: Falha ao armazenar cozinha %d no cache: %v", id, err)
		}
	}

	return cozinha, nil
}

func (s *CozinhaService) Save(cozinha *model.Cozinha) error {
	err := s.repo.Save(cozinha)
	if err != nil {
		return err
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateCozinha(cozinha.ID)
	}

	return nil
}

func (s *CozinhaService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Cozinha nao pode ser removida, pois esta em uso")
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateCozinha(id)
	}

	return nil
}
