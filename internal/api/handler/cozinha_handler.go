package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/assembler"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/domain/service"
	"github.com/yurisasc/algafood-go/pkg/pagination"
)

type CozinhaHandler struct {
	service *service.CozinhaService
}

func NewCozinhaHandler(service *service.CozinhaService) *CozinhaHandler {
	return &CozinhaHandler{service: service}
}

func (h *CozinhaHandler) Listar(c *gin.Context) {
	page := pagination.NewPageableFromContext(c)
	result, err := h.service.FindAll(page)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Convert to DTOs
	content := assembler.ToCozinhaModels(result.Content)
	response := pagination.Page[dto.CozinhaModel]{
		Content:          content,
		TotalElements:    result.TotalElements,
		TotalPages:       result.TotalPages,
		Size:             result.Size,
		Number:           result.Number,
		NumberOfElements: result.NumberOfElements,
		First:            result.First,
		Last:             result.Last,
		Empty:            result.Empty,
	}

	c.JSON(http.StatusOK, response)
}

func (h *CozinhaHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cozinhaId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	cozinha, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToCozinhaModel(cozinha))
}

func (h *CozinhaHandler) Adicionar(c *gin.Context) {
	var input dto.CozinhaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	cozinha := assembler.ToCozinhaEntity(&input)
	if err := h.service.Save(cozinha); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToCozinhaModel(cozinha))
}

func (h *CozinhaHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cozinhaId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.CozinhaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	cozinha, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	cozinha.Nome = input.Nome
	if err := h.service.Save(cozinha); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToCozinhaModel(cozinha))
}

func (h *CozinhaHandler) Remover(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("cozinhaId"), 10, 64)
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
