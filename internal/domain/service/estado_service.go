package service

import (
	"errors"
	"log"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type EstadoService struct {
	repo     repository.EstadoRepository
	cacheSvc *LocationCacheService
}

func NewEstadoService(repo repository.EstadoRepository, cacheSvc *LocationCacheService) *EstadoService {
	return &EstadoService{
		repo:     repo,
		cacheSvc: cacheSvc,
	}
}

func (s *EstadoService) FindAll() ([]model.Estado, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetAllEstados(); err == nil && cached != nil {
			return cached, nil
		}
	}

	estados, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetAllEstados(estados); err != nil {
			log.Printf("Aviso: Falha ao armazenar estados no cache: %v", err)
		}
	}

	return estados, nil
}

func (s *EstadoService) FindByID(id uint64) (*model.Estado, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetEstado(id); err == nil && cached != nil {
			return cached, nil
		}
	}

	estado, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewEstadoNaoEncontradoException(id)
		}
		return nil, err
	}

	// Armazena no cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetEstado(estado); err != nil {
			log.Printf("Aviso: Falha ao armazenar estado %d no cache: %v", id, err)
		}
	}

	return estado, nil
}

func (s *EstadoService) Save(estado *model.Estado) error {
	err := s.repo.Save(estado)
	if err != nil {
		return err
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateEstado(estado.ID)
	}

	return nil
}

func (s *EstadoService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Estado nao pode ser removido, pois esta em uso")
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateEstado(id)
	}

	return nil
}
