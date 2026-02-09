package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type EstadoService struct {
	repo repository.EstadoRepository
}

func NewEstadoService(repo repository.EstadoRepository) *EstadoService {
	return &EstadoService{repo: repo}
}

func (s *EstadoService) FindAll() ([]model.Estado, error) {
	return s.repo.FindAll()
}

func (s *EstadoService) FindByID(id uint64) (*model.Estado, error) {
	estado, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewEstadoNaoEncontradoException(id)
		}
		return nil, err
	}
	return estado, nil
}

func (s *EstadoService) Save(estado *model.Estado) error {
	return s.repo.Save(estado)
}

func (s *EstadoService) Delete(id uint64) error {
	// First check if exists
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Estado nao pode ser removido, pois esta em uso")
	}
	return nil
}
