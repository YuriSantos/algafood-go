package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// EstadoModel represents Estado output
type EstadoModel struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// CidadeModel represents Cidade output
type CidadeModel struct {
	ID     uint64      `json:"id"`
	Nome   string      `json:"nome"`
	Estado EstadoModel `json:"estado"`
}

// CozinhaModel represents Cozinha output
type CozinhaModel struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// FormaPagamentoModel represents FormaPagamento output
type FormaPagamentoModel struct {
	ID              uint64    `json:"id"`
	Descricao       string    `json:"descricao"`
	DataAtualizacao time.Time `json:"dataAtualizacao"`
}

// PermissaoModel represents Permissao output
type PermissaoModel struct {
	ID        uint64 `json:"id"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}

// GrupoModel represents Grupo output
type GrupoModel struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// UsuarioModel represents Usuario output
type UsuarioModel struct {
	ID           uint64    `json:"id"`
	Nome         string    `json:"nome"`
	Email        string    `json:"email"`
	DataCadastro time.Time `json:"dataCadastro"`
}

// RestauranteModel represents full Restaurante output
type RestauranteModel struct {
	ID              uint64          `json:"id"`
	Nome            string          `json:"nome"`
	TaxaFrete       decimal.Decimal `json:"taxaFrete"`
	Cozinha         CozinhaModel    `json:"cozinha"`
	Ativo           bool            `json:"ativo"`
	Aberto          bool            `json:"aberto"`
	Endereco        *EnderecoModel  `json:"endereco,omitempty"`
	DataCadastro    time.Time       `json:"dataCadastro"`
	DataAtualizacao time.Time       `json:"dataAtualizacao"`
}

// RestauranteResumoModel represents summary Restaurante output
type RestauranteResumoModel struct {
	ID        uint64          `json:"id"`
	Nome      string          `json:"nome"`
	TaxaFrete decimal.Decimal `json:"taxaFrete"`
	Cozinha   CozinhaModel    `json:"cozinha"`
	Ativo     bool            `json:"ativo"`
	Aberto    bool            `json:"aberto"`
}

// RestauranteApenasNomeModel represents minimal Restaurante output
type RestauranteApenasNomeModel struct {
	ID   uint64 `json:"id"`
	Nome string `json:"nome"`
}

// EnderecoModel represents address output
type EnderecoModel struct {
	CEP         string            `json:"cep"`
	Logradouro  string            `json:"logradouro"`
	Numero      string            `json:"numero"`
	Complemento string            `json:"complemento,omitempty"`
	Bairro      string            `json:"bairro"`
	Cidade      CidadeResumoModel `json:"cidade"`
}

// CidadeResumoModel represents summary Cidade output
type CidadeResumoModel struct {
	ID     uint64 `json:"id"`
	Nome   string `json:"nome"`
	Estado string `json:"estado"`
}

// ProdutoModel represents Produto output
type ProdutoModel struct {
	ID        uint64          `json:"id"`
	Nome      string          `json:"nome"`
	Descricao string          `json:"descricao"`
	Preco     decimal.Decimal `json:"preco"`
	Ativo     bool            `json:"ativo"`
}

// FotoProdutoModel represents FotoProduto output
type FotoProdutoModel struct {
	NomeArquivo string `json:"nomeArquivo"`
	Descricao   string `json:"descricao"`
	ContentType string `json:"contentType"`
	Tamanho     int64  `json:"tamanho"`
}

// PedidoModel represents full Pedido output
type PedidoModel struct {
	Codigo           string                     `json:"codigo"`
	Subtotal         decimal.Decimal            `json:"subtotal"`
	TaxaFrete        decimal.Decimal            `json:"taxaFrete"`
	ValorTotal       decimal.Decimal            `json:"valorTotal"`
	Status           string                     `json:"status"`
	DataCriacao      time.Time                  `json:"dataCriacao"`
	DataConfirmacao  *time.Time                 `json:"dataConfirmacao,omitempty"`
	DataCancelamento *time.Time                 `json:"dataCancelamento,omitempty"`
	DataEntrega      *time.Time                 `json:"dataEntrega,omitempty"`
	Restaurante      RestauranteApenasNomeModel `json:"restaurante"`
	Cliente          UsuarioModel               `json:"cliente"`
	FormaPagamento   FormaPagamentoModel        `json:"formaPagamento"`
	EnderecoEntrega  EnderecoModel              `json:"enderecoEntrega"`
	Itens            []ItemPedidoModel          `json:"itens"`
}

// PedidoResumoModel represents summary Pedido output
type PedidoResumoModel struct {
	Codigo      string                     `json:"codigo"`
	Subtotal    decimal.Decimal            `json:"subtotal"`
	TaxaFrete   decimal.Decimal            `json:"taxaFrete"`
	ValorTotal  decimal.Decimal            `json:"valorTotal"`
	Status      string                     `json:"status"`
	DataCriacao time.Time                  `json:"dataCriacao"`
	Restaurante RestauranteApenasNomeModel `json:"restaurante"`
	Cliente     UsuarioModel               `json:"cliente"`
}

// ItemPedidoModel represents ItemPedido output
type ItemPedidoModel struct {
	ProdutoID     uint64          `json:"produtoId"`
	ProdutoNome   string          `json:"produtoNome"`
	Quantidade    int             `json:"quantidade"`
	PrecoUnitario decimal.Decimal `json:"precoUnitario"`
	PrecoTotal    decimal.Decimal `json:"precoTotal"`
	Observacao    string          `json:"observacao,omitempty"`
}

// VendaDiariaModel represents daily sales output
type VendaDiariaModel struct {
	Data          string  `json:"data"`
	TotalVendas   int64   `json:"totalVendas"`
	TotalFaturado float64 `json:"totalFaturado"`
}
