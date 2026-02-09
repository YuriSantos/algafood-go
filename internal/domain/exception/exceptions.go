package exception

import "fmt"

// NegocioException represents a business logic error
type NegocioException struct {
	Message string
}

func (e *NegocioException) Error() string {
	return e.Message
}

func NewNegocioException(message string) *NegocioException {
	return &NegocioException{Message: message}
}

// EntidadeNaoEncontradaException represents a not found error
type EntidadeNaoEncontradaException struct {
	Message string
}

func (e *EntidadeNaoEncontradaException) Error() string {
	return e.Message
}

// EntidadeEmUsoException represents an entity in use error (cannot delete)
type EntidadeEmUsoException struct {
	Message string
}

func (e *EntidadeEmUsoException) Error() string {
	return e.Message
}

func NewEntidadeEmUsoException(message string) *EntidadeEmUsoException {
	return &EntidadeEmUsoException{Message: message}
}

// Specific entity not found exceptions

type EstadoNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewEstadoNaoEncontradoException(estadoID uint64) *EstadoNaoEncontradoException {
	return &EstadoNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de estado com codigo %d", estadoID),
		},
	}
}

type CidadeNaoEncontradaException struct {
	EntidadeNaoEncontradaException
}

func NewCidadeNaoEncontradaException(cidadeID uint64) *CidadeNaoEncontradaException {
	return &CidadeNaoEncontradaException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de cidade com codigo %d", cidadeID),
		},
	}
}

type CozinhaNaoEncontradaException struct {
	EntidadeNaoEncontradaException
}

func NewCozinhaNaoEncontradaException(cozinhaID uint64) *CozinhaNaoEncontradaException {
	return &CozinhaNaoEncontradaException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de cozinha com codigo %d", cozinhaID),
		},
	}
}

type RestauranteNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewRestauranteNaoEncontradoException(restauranteID uint64) *RestauranteNaoEncontradoException {
	return &RestauranteNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de restaurante com codigo %d", restauranteID),
		},
	}
}

type ProdutoNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewProdutoNaoEncontradoException(restauranteID, produtoID uint64) *ProdutoNaoEncontradoException {
	return &ProdutoNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de produto com codigo %d para o restaurante de codigo %d", produtoID, restauranteID),
		},
	}
}

type FormaPagamentoNaoEncontradaException struct {
	EntidadeNaoEncontradaException
}

func NewFormaPagamentoNaoEncontradaException(formaPagamentoID uint64) *FormaPagamentoNaoEncontradaException {
	return &FormaPagamentoNaoEncontradaException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de forma de pagamento com codigo %d", formaPagamentoID),
		},
	}
}

type UsuarioNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewUsuarioNaoEncontradoException(usuarioID uint64) *UsuarioNaoEncontradoException {
	return &UsuarioNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de usuario com codigo %d", usuarioID),
		},
	}
}

type GrupoNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewGrupoNaoEncontradoException(grupoID uint64) *GrupoNaoEncontradoException {
	return &GrupoNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de grupo com codigo %d", grupoID),
		},
	}
}

type PermissaoNaoEncontradaException struct {
	EntidadeNaoEncontradaException
}

func NewPermissaoNaoEncontradaException(permissaoID uint64) *PermissaoNaoEncontradaException {
	return &PermissaoNaoEncontradaException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de permissao com codigo %d", permissaoID),
		},
	}
}

type PedidoNaoEncontradoException struct {
	EntidadeNaoEncontradaException
}

func NewPedidoNaoEncontradoException(codigoPedido string) *PedidoNaoEncontradoException {
	return &PedidoNaoEncontradoException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um pedido com codigo %s", codigoPedido),
		},
	}
}

type FotoProdutoNaoEncontradaException struct {
	EntidadeNaoEncontradaException
}

func NewFotoProdutoNaoEncontradaException(restauranteID, produtoID uint64) *FotoProdutoNaoEncontradaException {
	return &FotoProdutoNaoEncontradaException{
		EntidadeNaoEncontradaException{
			Message: fmt.Sprintf("Nao existe um cadastro de foto do produto com codigo %d para o restaurante de codigo %d", produtoID, restauranteID),
		},
	}
}
