package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/assembler"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/domain/service"
)

type ProdutoHandler struct {
	service *service.ProdutoService
}

func NewProdutoHandler(service *service.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{service: service}
}

func (h *ProdutoHandler) Listar(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	incluirInativos := c.Query("incluirInativos") == "true"

	produtos, err := h.service.FindAllByRestaurante(restauranteID, incluirInativos)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToProdutoModels(produtos))
}

func (h *ProdutoHandler) Buscar(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	produtoID, _ := strconv.ParseUint(c.Param("produtoId"), 10, 64)

	produto, err := h.service.FindByID(restauranteID, produtoID)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToProdutoModel(produto))
}

func (h *ProdutoHandler) Adicionar(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)

	var input dto.ProdutoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	produto := assembler.ToProdutoEntity(&input)
	if err := h.service.Save(restauranteID, produto); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToProdutoModel(produto))
}

func (h *ProdutoHandler) Atualizar(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	produtoID, _ := strconv.ParseUint(c.Param("produtoId"), 10, 64)

	var input dto.ProdutoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	produto, err := h.service.FindByID(restauranteID, produtoID)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Update fields
	updated := assembler.ToProdutoEntity(&input)
	produto.Nome = updated.Nome
	produto.Descricao = updated.Descricao
	produto.Preco = updated.Preco
	produto.Ativo = updated.Ativo

	if err := h.service.Save(restauranteID, produto); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToProdutoModel(produto))
}
