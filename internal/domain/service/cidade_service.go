package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type CidadeService struct {
	repo      repository.CidadeRepository
	estadoSvc *EstadoService
}

func NewCidadeService(repo repository.CidadeRepository, estadoSvc *EstadoService) *CidadeService {
	return &CidadeService{
		repo:      repo,
		estadoSvc: estadoSvc,
	}
}

func (s *CidadeService) FindAll() ([]model.Cidade, error) {
	return s.repo.FindAll()
}

func (s *CidadeService) FindByID(id uint64) (*model.Cidade, error) {
	cidade, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewCidadeNaoEncontradaException(id)
		}
		return nil, err
	}
	return cidade, nil
}

func (s *CidadeService) Save(cidade *model.Cidade) error {
	// Validate estado exists
	estado, err := s.estadoSvc.FindByID(cidade.EstadoID)
	if err != nil {
		return err
	}
	cidade.Estado = *estado

	return s.repo.Save(cidade)
}

func (s *CidadeService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Cidade nao pode ser removida, pois esta em uso")
	}
	return nil
}
