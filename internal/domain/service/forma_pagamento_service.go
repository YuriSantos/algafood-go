package service

import (
	"errors"
	"log"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type FormaPagamentoService struct {
	repo     repository.FormaPagamentoRepository
	cacheSvc *BusinessCacheService
}

func NewFormaPagamentoService(repo repository.FormaPagamentoRepository, cacheSvc *BusinessCacheService) *FormaPagamentoService {
	return &FormaPagamentoService{
		repo:     repo,
		cacheSvc: cacheSvc,
	}
}

func (s *FormaPagamentoService) FindAll() ([]model.FormaPagamento, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetAllFormasPagamento(); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Cache miss - busca do banco
	fps, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Atualiza o cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetAllFormasPagamento(fps); err != nil {
			log.Printf("Aviso: Falha ao armazenar formas de pagamento no cache: %v", err)
		}
	}

	return fps, nil
}

func (s *FormaPagamentoService) FindByID(id uint64) (*model.FormaPagamento, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetFormaPagamento(id); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Cache miss - busca do banco
	fp, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewFormaPagamentoNaoEncontradaException(id)
		}
		return nil, err
	}

	// Atualiza o cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetFormaPagamento(fp); err != nil {
			log.Printf("Aviso: Falha ao armazenar forma de pagamento %d no cache: %v", id, err)
		}
	}

	return fp, nil
}

func (s *FormaPagamentoService) Save(fp *model.FormaPagamento) error {
	err := s.repo.Save(fp)
	if err != nil {
		return err
	}

	// Invalida o cache após salvar
	if s.cacheSvc != nil {
		if err := s.cacheSvc.InvalidateFormaPagamento(fp.ID); err != nil {
			log.Printf("Aviso: Falha ao invalidar cache da forma de pagamento %d: %v", fp.ID, err)
		}
	}

	return nil
}

func (s *FormaPagamentoService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Forma de pagamento nao pode ser removida, pois esta em uso")
	}

	// Invalida o cache após deletar
	if s.cacheSvc != nil {
		if err := s.cacheSvc.InvalidateFormaPagamento(id); err != nil {
			log.Printf("Aviso: Falha ao invalidar cache da forma de pagamento %d: %v", id, err)
		}
	}

	return nil
}
