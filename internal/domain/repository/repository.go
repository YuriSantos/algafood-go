package repository

import (
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/pkg/pagination"
)

// EstadoRepository interface for estado operations
type EstadoRepository interface {
	FindAll() ([]model.Estado, error)
	FindByID(id uint64) (*model.Estado, error)
	Save(estado *model.Estado) error
	Delete(id uint64) error
}

// CidadeRepository interface for cidade operations
type CidadeRepository interface {
	FindAll() ([]model.Cidade, error)
	FindByID(id uint64) (*model.Cidade, error)
	Save(cidade *model.Cidade) error
	Delete(id uint64) error
}

// CozinhaRepository interface for cozinha operations
type CozinhaRepository interface {
	FindAll(page *pagination.Pageable) (*pagination.Page[model.Cozinha], error)
	FindByID(id uint64) (*model.Cozinha, error)
	Save(cozinha *model.Cozinha) error
	Delete(id uint64) error
}

// FormaPagamentoRepository interface for forma_pagamento operations
type FormaPagamentoRepository interface {
	FindAll() ([]model.FormaPagamento, error)
	FindByID(id uint64) (*model.FormaPagamento, error)
	Save(formaPagamento *model.FormaPagamento) error
	Delete(id uint64) error
}

// PermissaoRepository interface for permissao operations
type PermissaoRepository interface {
	FindAll() ([]model.Permissao, error)
	FindByID(id uint64) (*model.Permissao, error)
}

// GrupoRepository interface for grupo operations
type GrupoRepository interface {
	FindAll() ([]model.Grupo, error)
	FindByID(id uint64) (*model.Grupo, error)
	Save(grupo *model.Grupo) error
	Delete(id uint64) error
	AddPermissao(grupoID, permissaoID uint64) error
	RemovePermissao(grupoID, permissaoID uint64) error
}

// UsuarioRepository interface for usuario operations
type UsuarioRepository interface {
	FindAll() ([]model.Usuario, error)
	FindByID(id uint64) (*model.Usuario, error)
	FindByEmail(email string) (*model.Usuario, error)
	Save(usuario *model.Usuario) error
	AddGrupo(usuarioID, grupoID uint64) error
	RemoveGrupo(usuarioID, grupoID uint64) error
}

// RestauranteRepository interface for restaurante operations
type RestauranteRepository interface {
	FindAll() ([]model.Restaurante, error)
	FindByID(id uint64) (*model.Restaurante, error)
	Save(restaurante *model.Restaurante) error
	AddFormaPagamento(restauranteID, formaPagamentoID uint64) error
	RemoveFormaPagamento(restauranteID, formaPagamentoID uint64) error
	AddResponsavel(restauranteID, usuarioID uint64) error
	RemoveResponsavel(restauranteID, usuarioID uint64) error
	ExistsResponsavel(restauranteID, usuarioID uint64) (bool, error)
}

// ProdutoRepository interface for produto operations
type ProdutoRepository interface {
	FindAllByRestaurante(restauranteID uint64, incluirInativos bool) ([]model.Produto, error)
	FindByID(restauranteID, produtoID uint64) (*model.Produto, error)
	Save(produto *model.Produto) error
}

// FotoProdutoRepository interface for foto_produto operations
type FotoProdutoRepository interface {
	FindByProdutoID(produtoID uint64) (*model.FotoProduto, error)
	Save(foto *model.FotoProduto) error
	Delete(produtoID uint64) error
}

// PedidoFilter for filtering orders
type PedidoFilter struct {
	ClienteID         *uint64
	RestauranteID     *uint64
	DataCriacaoInicio *string
	DataCriacaoFim    *string
	Status            *model.StatusPedido
}

// PedidoRepository interface for pedido operations
type PedidoRepository interface {
	FindAll(filter *PedidoFilter, page *pagination.Pageable) (*pagination.Page[model.Pedido], error)
	FindByCodigo(codigo string) (*model.Pedido, error)
	Save(pedido *model.Pedido) error
	IsPedidoGerenciadoPor(codigoPedido string, usuarioID uint64) (bool, error)
}

// VendaDiaria represents daily sales statistics
type VendaDiaria struct {
	Data          string  `json:"data"`
	TotalVendas   int64   `json:"totalVendas"`
	TotalFaturado float64 `json:"totalFaturado"`
}

// VendaDiariaFilter for filtering sales report
type VendaDiariaFilter struct {
	RestauranteID     *uint64
	DataCriacaoInicio *string
	DataCriacaoFim    *string
}

// VendaQueryRepository interface for sales queries
type VendaQueryRepository interface {
	ConsultarVendasDiarias(filter *VendaDiariaFilter, timeOffset string) ([]VendaDiaria, error)
}
