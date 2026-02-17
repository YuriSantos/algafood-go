package service

import (
	"errors"
	"log"

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
	cacheSvc          *BusinessCacheService
}

func NewRestauranteService(
	repo repository.RestauranteRepository,
	cozinhaSvc *CozinhaService,
	cidadeSvc *CidadeService,
	formaPagamentoSvc *FormaPagamentoService,
	usuarioSvc *UsuarioService,
	cacheSvc *BusinessCacheService,
) *RestauranteService {
	return &RestauranteService{
		repo:              repo,
		cozinhaSvc:        cozinhaSvc,
		cidadeSvc:         cidadeSvc,
		formaPagamentoSvc: formaPagamentoSvc,
		usuarioSvc:        usuarioSvc,
		cacheSvc:          cacheSvc,
	}
}

func (s *RestauranteService) FindAll() ([]model.Restaurante, error) {
	return s.repo.FindAll()
}

func (s *RestauranteService) FindByID(id uint64) (*model.Restaurante, error) {
	// Tenta obter do cache primeiro
	if s.cacheSvc != nil {
		if cached, err := s.cacheSvc.GetRestaurante(id); err == nil && cached != nil {
			// Converte para model e popula relacionamentos do cache
			restaurante := cached.ToModel()

			// Popula cozinha do cache
			if cached.CozinhaID > 0 {
				if cozinha, err := s.cozinhaSvc.FindByID(cached.CozinhaID); err == nil {
					restaurante.Cozinha = *cozinha
				}
			}

			// Popula cidade do endereço do cache
			if cached.EnderecoCidadeID > 0 {
				if cidade, err := s.cidadeSvc.FindByID(cached.EnderecoCidadeID); err == nil {
					restaurante.Endereco.Cidade = *cidade
					restaurante.Endereco.CidadeID = cidade.ID
				}
			}

			// Popula formas de pagamento do cache
			for _, fpID := range cached.FormasPagamentoIDs {
				if fp, err := s.formaPagamentoSvc.FindByID(fpID); err == nil {
					restaurante.FormasPagamento = append(restaurante.FormasPagamento, *fp)
				}
			}

			// Popula responsáveis do cache
			for _, respID := range cached.ResponsaveisIDs {
				if resp, err := s.usuarioSvc.FindByID(respID); err == nil {
					restaurante.Responsaveis = append(restaurante.Responsaveis, *resp)
				}
			}

			return restaurante, nil
		}
	}

	// Cache miss - busca do banco
	restaurante, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewRestauranteNaoEncontradoException(id)
		}
		return nil, err
	}

	// Atualiza o cache
	if s.cacheSvc != nil {
		if err := s.cacheSvc.SetRestaurante(restaurante); err != nil {
			log.Printf("Aviso: Falha ao armazenar restaurante %d no cache: %v", id, err)
		}
	}

	return restaurante, nil
}

func (s *RestauranteService) Save(restaurante *model.Restaurante) error {
	// Validate cozinha exists (usa cache)
	cozinha, err := s.cozinhaSvc.FindByID(restaurante.CozinhaID)
	if err != nil {
		return err
	}
	restaurante.Cozinha = *cozinha

	// Validate cidade if endereco is provided (usa cache)
	if restaurante.Endereco.CidadeID != 0 {
		cidade, err := s.cidadeSvc.FindByID(restaurante.Endereco.CidadeID)
		if err != nil {
			return err
		}
		restaurante.Endereco.Cidade = *cidade
	}

	err = s.repo.Save(restaurante)
	if err != nil {
		return err
	}

	// Invalida o cache após salvar
	if s.cacheSvc != nil {
		if err := s.cacheSvc.InvalidateRestaurante(restaurante.ID); err != nil {
			log.Printf("Aviso: Falha ao invalidar cache do restaurante %d: %v", restaurante.ID, err)
		}
	}

	return nil
}

func (s *RestauranteService) Ativar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Ativar()
	err = s.repo.Save(restaurante)
	if err != nil {
		return err
	}
	s.invalidateCache(id)
	return nil
}

func (s *RestauranteService) Inativar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Inativar()
	err = s.repo.Save(restaurante)
	if err != nil {
		return err
	}
	s.invalidateCache(id)
	return nil
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
	err = s.repo.Save(restaurante)
	if err != nil {
		return err
	}
	s.invalidateCache(id)
	return nil
}

func (s *RestauranteService) Fechar(id uint64) error {
	restaurante, err := s.FindByID(id)
	if err != nil {
		return err
	}
	restaurante.Fechar()
	err = s.repo.Save(restaurante)
	if err != nil {
		return err
	}
	s.invalidateCache(id)
	return nil
}

func (s *RestauranteService) AssociarFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.formaPagamentoSvc.FindByID(formaPagamentoID); err != nil {
		return err
	}
	err := s.repo.AddFormaPagamento(restauranteID, formaPagamentoID)
	if err != nil {
		return err
	}
	s.invalidateCache(restauranteID)
	return nil
}

func (s *RestauranteService) DesassociarFormaPagamento(restauranteID, formaPagamentoID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.formaPagamentoSvc.FindByID(formaPagamentoID); err != nil {
		return err
	}
	err := s.repo.RemoveFormaPagamento(restauranteID, formaPagamentoID)
	if err != nil {
		return err
	}
	s.invalidateCache(restauranteID)
	return nil
}

func (s *RestauranteService) AssociarResponsavel(restauranteID, usuarioID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.usuarioSvc.FindByID(usuarioID); err != nil {
		return err
	}
	err := s.repo.AddResponsavel(restauranteID, usuarioID)
	if err != nil {
		return err
	}
	s.invalidateCache(restauranteID)
	return nil
}

func (s *RestauranteService) DesassociarResponsavel(restauranteID, usuarioID uint64) error {
	if _, err := s.FindByID(restauranteID); err != nil {
		return err
	}
	if _, err := s.usuarioSvc.FindByID(usuarioID); err != nil {
		return err
	}
	err := s.repo.RemoveResponsavel(restauranteID, usuarioID)
	if err != nil {
		return err
	}
	s.invalidateCache(restauranteID)
	return nil
}

// invalidateCache invalida o cache de um restaurante
func (s *RestauranteService) invalidateCache(id uint64) {
	if s.cacheSvc != nil {
		if err := s.cacheSvc.InvalidateRestaurante(id); err != nil {
			log.Printf("Aviso: Falha ao invalidar cache do restaurante %d: %v", id, err)
		}
	}
}
