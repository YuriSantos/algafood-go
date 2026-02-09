package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type GrupoService struct {
	repo         repository.GrupoRepository
	permissaoSvc *PermissaoService
}

func NewGrupoService(repo repository.GrupoRepository, permissaoSvc *PermissaoService) *GrupoService {
	return &GrupoService{
		repo:         repo,
		permissaoSvc: permissaoSvc,
	}
}

func (s *GrupoService) FindAll() ([]model.Grupo, error) {
	return s.repo.FindAll()
}

func (s *GrupoService) FindByID(id uint64) (*model.Grupo, error) {
	grupo, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewGrupoNaoEncontradoException(id)
		}
		return nil, err
	}
	return grupo, nil
}

func (s *GrupoService) Save(grupo *model.Grupo) error {
	return s.repo.Save(grupo)
}

func (s *GrupoService) Delete(id uint64) error {
	if _, err := s.FindByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return exception.NewEntidadeEmUsoException("Grupo nao pode ser removido, pois esta em uso")
	}
	return nil
}

func (s *GrupoService) AssociarPermissao(grupoID, permissaoID uint64) error {
	if _, err := s.FindByID(grupoID); err != nil {
		return err
	}
	if _, err := s.permissaoSvc.FindByID(permissaoID); err != nil {
		return err
	}
	return s.repo.AddPermissao(grupoID, permissaoID)
}

func (s *GrupoService) DesassociarPermissao(grupoID, permissaoID uint64) error {
	if _, err := s.FindByID(grupoID); err != nil {
		return err
	}
	if _, err := s.permissaoSvc.FindByID(permissaoID); err != nil {
		return err
	}
	return s.repo.RemovePermissao(grupoID, permissaoID)
}
