package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type FormaPagamentoService struct {
	repo repository.FormaPagamentoRepository
}

func NewFormaPagamentoService(repo repository.FormaPagamentoRepository) *FormaPagamentoService {
	return &FormaPagamentoService{repo: repo}
}

func (s *FormaPagamentoService) FindAll() ([]model.FormaPagamento, error) {
	return s.repo.FindAll()
}

func (s *FormaPagamentoService) FindByID(id uint64) (*model.FormaPagamento, error) {
	fp, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewFormaPagamentoNaoEncontradaException(id)
		}
		return nil, err
	}
	return fp, nil
}

func (s *FormaPagamentoService) Save(fp *model.FormaPagamento) error {
	return s.repo.Save(fp)
}

func (s *FormaPagamentoService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Forma de pagamento nao pode ser removida, pois esta em uso")
	}
	return nil
}
