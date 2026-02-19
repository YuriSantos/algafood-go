package service

import (
	"errors"
	"log"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type CidadeService struct {
	repo      repository.CidadeRepository
	estadoSvc *EstadoService
	cacheSvc  *LocationCacheService
}

func NewCidadeService(repo repository.CidadeRepository, estadoSvc *EstadoService, cacheSvc *LocationCacheService) *CidadeService {
	return &CidadeService{
		repo:      repo,
		estadoSvc: estadoSvc,
		cacheSvc:  cacheSvc,
	}
}

func (s *CidadeService) FindAll() ([]model.Cidade, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetAllCidades(); err == nil && cached != nil {
			return cached, nil
		}
	}

	cidades, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Armazena no cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetAllCidades(cidades); err != nil {
			log.Printf("Aviso: Falha ao armazenar cidades no cache: %v", err)
		}
	}

	return cidades, nil
}

func (s *CidadeService) FindByID(id uint64) (*model.Cidade, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetCidade(id); err == nil && cached != nil {
			return cached, nil
		}
	}

	cidade, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewCidadeNaoEncontradaException(id)
		}
		return nil, err
	}

	// Armazena no cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetCidade(cidade); err != nil {
			log.Printf("Aviso: Falha ao armazenar cidade %d no cache: %v", id, err)
		}
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

	err = s.repo.Save(cidade)
	if err != nil {
		return err
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateCidade(cidade)
	}

	return nil
}

func (s *CidadeService) Delete(id uint64) error {
	cidade, err := s.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Cidade nao pode ser removida, pois esta em uso")
	}

	// Invalida cache
	if s.cacheSvc != nil {
		s.cacheSvc.InvalidateCidade(cidade)
	}

	return nil
}
