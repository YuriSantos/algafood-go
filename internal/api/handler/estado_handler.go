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

type EstadoHandler struct {
	service *service.EstadoService
}

func NewEstadoHandler(service *service.EstadoService) *EstadoHandler {
	return &EstadoHandler{service: service}
}

func (h *EstadoHandler) Listar(c *gin.Context) {
	estados, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToEstadoModels(estados))
}

func (h *EstadoHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("estadoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	estado, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToEstadoModel(estado))
}

func (h *EstadoHandler) Adicionar(c *gin.Context) {
	var input dto.EstadoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	estado := assembler.ToEstadoEntity(&input)
	if err := h.service.Save(estado); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToEstadoModel(estado))
}

func (h *EstadoHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("estadoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.EstadoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	estado, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	estado.Nome = input.Nome
	if err := h.service.Save(estado); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToEstadoModel(estado))
}

func (h *EstadoHandler) Remover(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("estadoId"), 10, 64)
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
