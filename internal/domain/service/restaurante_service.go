package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type RestauranteService struct {
	repo              repository.RestauranteRepository
	cozinhaSvc        *CozinhaService
	cidadeSvc         *CidadeService
	formaPagamentoSvc *FormaPagamentoService
	usuarioSvc        *UsuarioService
}

func NewRestauranteService(
	repo repository.RestauranteRepository,
	cozinhaSvc *CozinhaService,
	cidadeSvc *CidadeService,
	formaPagamentoSvc *FormaPagamentoService,
	usuarioSvc *UsuarioService,
) *RestauranteService {
	return &RestauranteService{
		repo:              repo,
		cozinhaSvc:        cozinhaSvc,
		cidadeSvc:         cidadeSvc,
		formaPagamentoSvc: formaPagamentoSvc,
		usuarioSvc:        usuarioSvc,
	}
}

func (s *RestauranteService) FindAll() ([]model.Restaurante, error) {
	return s.repo.FindAll()
}

func (s *RestauranteService) FindByID(id uint64) (*model.Restaurante, error) {
	restaurante, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewRestauranteNaoEncontradoException(id)
		}
		return nil, err
	}
	return restaurante, nil
}

func (s *RestauranteService) Save(restaurante *model.Restaurante) error {
	// Validate cozinha exists
	cozinha, err := s.cozinhaSvc.FindByID(restaurante.CozinhaID)
	if err != nil {
		return err
	}
	restaurante.Cozinha = *cozinha

	// Validate cidade if endereco is provided
	if restaurante.Endereco.CidadeID != 0 {
		cidade, err := s.cidadeSvc.FindByID(restaurante.Endereco.CidadeID)
		if err != nil {
			return err
		}
		restaurante.Endereco.Cidade = *cidade
	}

	return s.repo.Save(restaurante)
}

func (s *RestauranteService) Ativar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Ativar()
	return s.repo.Save(restaurante)
}

func (s *RestauranteService) Inativar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Inativar()
	return s.repo.Save(restaurante)
}

func (s *RestauranteService) AtivarEmMassa(ids []uint64) error {
	for _, id := range ids {
		if err := s.Ativar(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *RestauranteService) InativarEmMassa(ids []uint64) error {
	for _, id := range ids {
		if err := s.Inativar(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *RestauranteService) Abrir(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Abrir()
	return s.repo.Save(restaurante)
}

func (s *RestauranteService) Fechar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Fechar()
	return s.repo.Save(restaurante)
}

func (s *RestauranteService) AssociarFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.formaPagamentoSvc.FindByID(formaPagamentoID); err != nil {
		return err
	}
	return s.repo.AddFormaPagamento(restauranteID, formaPagamentoID)
}

func (s *RestauranteService) DesassociarFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.formaPagamentoSvc.FindByID(formaPagamentoID); err != nil {
		return err
	}
	return s.repo.RemoveFormaPagamento(restauranteID, formaPagamentoID)
}

func (s *RestauranteService) AssociarResponsavel(restauranteID, usuarioID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.usuarioSvc.FindByID(usuarioID); err != nil {
		return err
	}
	return s.repo.AddResponsavel(restauranteID, usuarioID)
}

func (s *RestauranteService) DesassociarResponsavel(restauranteID, usuarioID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.usuarioSvc.FindByID(usuarioID); err != nil {
		return err
	}
	return s.repo.RemoveResponsavel(restauranteID, usuarioID)
}
