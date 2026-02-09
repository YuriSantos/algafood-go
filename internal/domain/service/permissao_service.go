package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type PermissaoService struct {
	repo repository.PermissaoRepository
}

func NewPermissaoService(repo repository.PermissaoRepository) *PermissaoService {
	return &PermissaoService{repo: repo}
}

func (s *PermissaoService) FindAll() ([]model.Permissao, error) {
	return s.repo.FindAll()
}

func (s *PermissaoService) FindByID(id uint64) (*model.Permissao, error) {
	permissao, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewPermissaoNaoEncontradaException(id)
		}
		return nil, err
	}
	return permissao, nil
}
