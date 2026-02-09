package dto

// EstadoInput represents input for creating/updating Estado
type EstadoInput struct {
	Nome string `json:"nome" binding:"required,min=2,max=80"`
}

// CidadeInput represents input for creating/updating Cidade
type CidadeInput struct {
	Nome   string        `json:"nome" binding:"required,min=2,max=80"`
	Estado EstadoIDInput `json:"estado" binding:"required"`
}

// EstadoIDInput represents Estado ID reference
type EstadoIDInput struct {
	ID uint64 `json:"id" binding:"required"`
}

// CozinhaInput represents input for creating/updating Cozinha
type CozinhaInput struct {
	Nome string `json:"nome" binding:"required,min=2,max=60"`
}

// FormaPagamentoInput represents input for creating/updating FormaPagamento
type FormaPagamentoInput struct {
	Descricao string `json:"descricao" binding:"required,min=2,max=60"`
}

// GrupoInput represents input for creating/updating Grupo
type GrupoInput struct {
	Nome string `json:"nome" binding:"required,min=2,max=60"`
}

// UsuarioInput represents input for updating Usuario (without password)
type UsuarioInput struct {
	Nome  string `json:"nome" binding:"required,min=2,max=80"`
	Email string `json:"email" binding:"required,email,max=255"`
}

// UsuarioComSenhaInput represents input for creating Usuario (with password)
type UsuarioComSenhaInput struct {
	Nome  string `json:"nome" binding:"required,min=2,max=80"`
	Email string `json:"email" binding:"required,email,max=255"`
	Senha string `json:"senha" binding:"required,min=6"`
}

// SenhaInput represents input for changing password
type SenhaInput struct {
	SenhaAtual string `json:"senhaAtual" binding:"required"`
	NovaSenha  string `json:"novaSenha" binding:"required,min=6"`
}

// RestauranteInput represents input for creating/updating Restaurante
type RestauranteInput struct {
	Nome      string         `json:"nome" binding:"required,min=2,max=80"`
	TaxaFrete float64        `json:"taxaFrete" binding:"required,gte=0"`
	Cozinha   CozinhaIDInput `json:"cozinha" binding:"required"`
	Endereco  *EnderecoInput `json:"endereco"`
}

// CozinhaIDInput represents Cozinha ID reference
type CozinhaIDInput struct {
	ID uint64 `json:"id" binding:"required"`
}

// EnderecoInput represents address input
type EnderecoInput struct {
	CEP         string        `json:"cep" binding:"required,max=9"`
	Logradouro  string        `json:"logradouro" binding:"required,max=100"`
	Numero      string        `json:"numero" binding:"required,max=20"`
	Complemento string        `json:"complemento" binding:"max=60"`
	Bairro      string        `json:"bairro" binding:"required,max=60"`
	Cidade      CidadeIDInput `json:"cidade" binding:"required"`
}

// CidadeIDInput represents Cidade ID reference
type CidadeIDInput struct {
	ID uint64 `json:"id" binding:"required"`
}

// ProdutoInput represents input for creating/updating Produto
type ProdutoInput struct {
	Nome      string  `json:"nome" binding:"required,min=2,max=80"`
	Descricao string  `json:"descricao" binding:"max=255"`
	Preco     float64 `json:"preco" binding:"required,gt=0"`
	Ativo     bool    `json:"ativo"`
}

// PedidoInput represents input for creating Pedido
type PedidoInput struct {
	Restaurante     RestauranteIDInput    `json:"restaurante" binding:"required"`
	FormaPagamento  FormaPagamentoIDInput `json:"formaPagamento" binding:"required"`
	EnderecoEntrega EnderecoInput         `json:"enderecoEntrega" binding:"required"`
	Itens           []ItemPedidoInput     `json:"itens" binding:"required,min=1,dive"`
}

// RestauranteIDInput represents Restaurante ID reference
type RestauranteIDInput struct {
	ID uint64 `json:"id" binding:"required"`
}

// FormaPagamentoIDInput represents FormaPagamento ID reference
type FormaPagamentoIDInput struct {
	ID uint64 `json:"id" binding:"required"`
}

// ItemPedidoInput represents input for Pedido items
type ItemPedidoInput struct {
	ProdutoID  uint64 `json:"produtoId" binding:"required"`
	Quantidade int    `json:"quantidade" binding:"required,min=1"`
	Observacao string `json:"observacao" binding:"max=255"`
}

// FotoProdutoInput represents input for uploading product photo
type FotoProdutoInput struct {
	Descricao string `form:"descricao" binding:"max=150"`
}

// AtivacaoRestauranteInput represents input for bulk activation
type AtivacaoRestauranteInput struct {
	IDs []uint64 `json:"restauranteIds" binding:"required,min=1"`
}
