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

type CidadeHandler struct {
	service *service.CidadeService
}

func NewCidadeHandler(service *service.CidadeService) *CidadeHandler {
	return &CidadeHandler{service: service}
}

func (h *CidadeHandler) Listar(c *gin.Context) {
	cidades, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToCidadeModels(cidades))
}

func (h *CidadeHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cidadeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	cidade, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToCidadeModel(cidade))
}

func (h *CidadeHandler) Adicionar(c *gin.Context) {
	var input dto.CidadeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	cidade := assembler.ToCidadeEntity(&input)
	if err := h.service.Save(cidade); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Reload to get Estado
	cidade, _ = h.service.FindByID(cidade.ID)
	c.JSON(http.StatusCreated, assembler.ToCidadeModel(cidade))
}

func (h *CidadeHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cidadeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.CidadeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	cidade, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	cidade.Nome = input.Nome
	cidade.EstadoID = input.Estado.ID

	if err := h.service.Save(cidade); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Reload to get Estado
	cidade, _ = h.service.FindByID(cidade.ID)
	c.JSON(http.StatusOK, assembler.ToCidadeModel(cidade))
}

func (h *CidadeHandler) Remover(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cidadeId"), 10, 64)
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
