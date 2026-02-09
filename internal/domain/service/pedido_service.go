package service

import (
	"errors"

	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type PedidoService struct {
	repo              repository.PedidoRepository
	restauranteSvc    *RestauranteService
	cidadeSvc         *CidadeService
	usuarioSvc        *UsuarioService
	produtoSvc        *ProdutoService
	formaPagamentoSvc *FormaPagamentoService
}

func NewPedidoService(
	repo repository.PedidoRepository,
	restauranteSvc *RestauranteService,
	cidadeSvc *CidadeService,
	usuarioSvc *UsuarioService,
	produtoSvc *ProdutoService,
	formaPagamentoSvc *FormaPagamentoService,
) *PedidoService {
	return &PedidoService{
		repo:              repo,
		restauranteSvc:    restauranteSvc,
		cidadeSvc:         cidadeSvc,
		usuarioSvc:        usuarioSvc,
		produtoSvc:        produtoSvc,
		formaPagamentoSvc: formaPagamentoSvc,
	}
}

func (s *PedidoService) Pesquisar(filter *repository.PedidoFilter, page *pagination.Pageable) (*pagination.Page[model.Pedido], error) {
	return s.repo.FindAll(filter, page)
}

func (s *PedidoService) FindByCodigo(codigo string) (*model.Pedido, error) {
	pedido, err := s.repo.FindByCodigo(codigo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewPedidoNaoEncontradoException(codigo)
		}
		return nil, err
	}
	return pedido, nil
}

func (s *PedidoService) Emitir(pedido *model.Pedido) error {
	// Validate restaurante
	restaurante, err := s.restauranteSvc.FindByID(pedido.RestauranteID)
	if err != nil {
		return err
	}
	pedido.Restaurante = *restaurante

	// Validate forma pagamento
	formaPagamento, err := s.formaPagamentoSvc.FindByID(pedido.FormaPagamentoID)
	if err != nil {
		return err
	}

	// Check if restaurante accepts this forma pagamento
	if restaurante.NaoAceitaFormaPagamento(*formaPagamento) {
		return exception.NewNegocioException("Forma de pagamento nao aceita por esse restaurante")
	}
	pedido.FormaPagamento = *formaPagamento

	// Validate cliente
	cliente, err := s.usuarioSvc.FindByID(pedido.ClienteID)
	if err != nil {
		return err
	}
	pedido.Cliente = *cliente

	// Validate cidade
	if pedido.EnderecoEntrega.CidadeID != 0 {
		cidade, err := s.cidadeSvc.FindByID(pedido.EnderecoEntrega.CidadeID)
		if err != nil {
			return err
		}
		pedido.EnderecoEntrega.Cidade = *cidade
	}

	// Validate and set items
	for i := range pedido.Itens {
		item := &pedido.Itens[i]
		produto, err := s.produtoSvc.FindByID(restaurante.ID, item.ProdutoID)
		if err != nil {
			return err
		}
		item.Produto = *produto
		item.PrecoUnitario = produto.Preco
		item.CalcularPrecoTotal()
	}

	// Set freight and calculate total
	pedido.TaxaFrete = restaurante.TaxaFrete
	pedido.CalcularValorTotal()
	pedido.BeforeCreate()

	return s.repo.Save(pedido)
}
