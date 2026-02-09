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

type FormaPagamentoHandler struct {
	service *service.FormaPagamentoService
}

func NewFormaPagamentoHandler(service *service.FormaPagamentoService) *FormaPagamentoHandler {
	return &FormaPagamentoHandler{service: service}
}

func (h *FormaPagamentoHandler) Listar(c *gin.Context) {
	formasPagamento, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToFormaPagamentoModels(formasPagamento))
}

func (h *FormaPagamentoHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("formaPagamentoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	fp, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToFormaPagamentoModel(fp))
}

func (h *FormaPagamentoHandler) Adicionar(c *gin.Context) {
	var input dto.FormaPagamentoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	fp := assembler.ToFormaPagamentoEntity(&input)
	if err := h.service.Save(fp); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToFormaPagamentoModel(fp))
}

func (h *FormaPagamentoHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("formaPagamentoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.FormaPagamentoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	fp, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	fp.Descricao = input.Descricao
	if err := h.service.Save(fp); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToFormaPagamentoModel(fp))
}

func (h *FormaPagamentoHandler) Remover(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("formaPagamentoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
