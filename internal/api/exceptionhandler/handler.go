package exceptionhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/domain/exception"
)

const (
	MSG_ERRO_GENERICA_USUARIO_FINAL = "Ocorreu um erro interno inesperado no sistema. Tente novamente e se o problema persistir, entre em contato com o administrador do sistema."
	MSG_DADOS_INVALIDOS             = "Um ou mais campos estao invalidos. Faca o preenchimento correto e tente novamente."
)

// HandleError handles domain exceptions and returns appropriate HTTP response
func HandleError(c *gin.Context, err error) {
	var entidadeNaoEncontrada *exception.EntidadeNaoEncontradaException
	var entidadeEmUso *exception.EntidadeEmUsoException
	var negocioException *exception.NegocioException
	var authenticationException *exception.AuthenticationException

	// Check for specific not found exceptions
	var estadoNaoEncontrado *exception.EstadoNaoEncontradoException
	var cidadeNaoEncontrada *exception.CidadeNaoEncontradaException
	var cozinhaNaoEncontrada *exception.CozinhaNaoEncontradaException
	var restauranteNaoEncontrado *exception.RestauranteNaoEncontradoException
	var produtoNaoEncontrado *exception.ProdutoNaoEncontradoException
	var formaPagamentoNaoEncontrada *exception.FormaPagamentoNaoEncontradaException
	var usuarioNaoEncontrado *exception.UsuarioNaoEncontradoException
	var grupoNaoEncontrado *exception.GrupoNaoEncontradoException
	var permissaoNaoEncontrada *exception.PermissaoNaoEncontradaException
	var pedidoNaoEncontrado *exception.PedidoNaoEncontradoException
	var fotoProdutoNaoEncontrada *exception.FotoProdutoNaoEncontradaException

	switch {
	case errors.As(err, &authenticationException):
		handleUnauthorized(c, authenticationException.Message)
	case errors.As(err, &estadoNaoEncontrado):
		handleNotFound(c, estadoNaoEncontrado.Message)
	case errors.As(err, &cidadeNaoEncontrada):
		handleNotFound(c, cidadeNaoEncontrada.Message)
	case errors.As(err, &cozinhaNaoEncontrada):
		handleNotFound(c, cozinhaNaoEncontrada.Message)
	case errors.As(err, &restauranteNaoEncontrado):
		handleNotFound(c, restauranteNaoEncontrado.Message)
	case errors.As(err, &produtoNaoEncontrado):
		handleNotFound(c, produtoNaoEncontrado.Message)
	case errors.As(err, &formaPagamentoNaoEncontrada):
		handleNotFound(c, formaPagamentoNaoEncontrada.Message)
	case errors.As(err, &usuarioNaoEncontrado):
		handleNotFound(c, usuarioNaoEncontrado.Message)
	case errors.As(err, &grupoNaoEncontrado):
		handleNotFound(c, grupoNaoEncontrado.Message)
	case errors.As(err, &permissaoNaoEncontrada):
		handleNotFound(c, permissaoNaoEncontrada.Message)
	case errors.As(err, &pedidoNaoEncontrado):
		handleNotFound(c, pedidoNaoEncontrado.Message)
	case errors.As(err, &fotoProdutoNaoEncontrada):
		handleNotFound(c, fotoProdutoNaoEncontrada.Message)
	case errors.As(err, &entidadeNaoEncontrada):
		handleNotFound(c, entidadeNaoEncontrada.Message)
	case errors.As(err, &entidadeEmUso):
		handleConflict(c, entidadeEmUso.Message)
	case errors.As(err, &negocioException):
		handleBadRequest(c, negocioException.Message)
	default:
		handleInternalError(c, err)
	}
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		objects := make([]dto.ObjectError, 0, len(validationErrors))
		for _, fieldError := range validationErrors {
			objects = append(objects, dto.ObjectError{
				Name:        fieldError.Field(),
				UserMessage: getValidationMessage(fieldError),
			})
		}

		problem := dto.NewProblemWithObjects(
			http.StatusBadRequest,
			dto.ProblemTypeInvalidData,
			MSG_DADOS_INVALIDOS,
			MSG_DADOS_INVALIDOS,
			objects,
		)
		c.JSON(http.StatusBadRequest, problem)
		return
	}

	// Handle JSON binding errors
	problem := dto.NewProblem(
		http.StatusBadRequest,
		dto.ProblemTypeInvalidMessage,
		"O corpo da requisicao esta invalido. Verifique erro de sintaxe.",
		"O corpo da requisicao esta invalido. Verifique erro de sintaxe.",
	)
	c.JSON(http.StatusBadRequest, problem)
}

func handleNotFound(c *gin.Context, message string) {
	problem := dto.NewProblem(
		http.StatusNotFound,
		dto.ProblemTypeResourceNotFound,
		message,
		message,
	)
	c.JSON(http.StatusNotFound, problem)
}

func handleConflict(c *gin.Context, message string) {
	problem := dto.NewProblem(
		http.StatusConflict,
		dto.ProblemTypeEntityInUse,
		message,
		message,
	)
	c.JSON(http.StatusConflict, problem)
}

func handleBadRequest(c *gin.Context, message string) {
	problem := dto.NewProblem(
		http.StatusBadRequest,
		dto.ProblemTypeBusinessError,
		message,
		message,
	)
	c.JSON(http.StatusBadRequest, problem)
}

func handleInternalError(c *gin.Context, err error) {
	problem := dto.NewProblem(
		http.StatusInternalServerError,
		dto.ProblemTypeSystemError,
		err.Error(),
		MSG_ERRO_GENERICA_USUARIO_FINAL,
	)
	c.JSON(http.StatusInternalServerError, problem)
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " e obrigatorio"
	case "email":
		return fe.Field() + " deve ser um e-mail valido"
	case "min":
		return fe.Field() + " deve ter no minimo " + fe.Param() + " caracteres"
	case "max":
		return fe.Field() + " deve ter no maximo " + fe.Param() + " caracteres"
	case "gte":
		return fe.Field() + " deve ser maior ou igual a " + fe.Param()
	case "gt":
		return fe.Field() + " deve ser maior que " + fe.Param()
	case "lte":
		return fe.Field() + " deve ser menor ou igual a " + fe.Param()
	case "lt":
		return fe.Field() + " deve ser menor que " + fe.Param()
	default:
		return fe.Field() + " esta invalido"
	}
}

// HandleAccessDenied handles authorization errors
func HandleAccessDenied(c *gin.Context) {
	problem := dto.NewProblem(
		http.StatusForbidden,
		dto.ProblemTypeAccessDenied,
		"Voce nao possui permissao para executar essa operacao.",
		"Voce nao possui permissao para executar essa operacao.",
	)
	c.JSON(http.StatusForbidden, problem)
}

// HandleUnauthorized handles authentication errors
func HandleUnauthorized(c *gin.Context) {
	problem := dto.NewProblem(
		http.StatusUnauthorized,
		dto.ProblemTypeAccessDenied,
		"Autenticacao necessaria para acessar esse recurso.",
		"Autenticacao necessaria para acessar esse recurso.",
	)
	c.JSON(http.StatusUnauthorized, problem)
}

func handleUnauthorized(c *gin.Context, message string) {
	problem := dto.NewProblem(
		http.StatusUnauthorized,
		dto.ProblemTypeInvalidCredentials,
		message,
		message,
	)
	c.JSON(http.StatusUnauthorized, problem)
}
