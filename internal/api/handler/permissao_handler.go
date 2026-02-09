package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/assembler"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/domain/service"
)

type PermissaoHandler struct {
	service *service.PermissaoService
}

func NewPermissaoHandler(service *service.PermissaoService) *PermissaoHandler {
	return &PermissaoHandler{service: service}
}

func (h *PermissaoHandler) Listar(c *gin.Context) {
	permissoes, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToPermissaoModels(permissoes))
}

func (h *PermissaoHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("permissaoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	permissao, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToPermissaoModel(permissao))
}
