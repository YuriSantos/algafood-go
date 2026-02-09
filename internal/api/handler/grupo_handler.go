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

type GrupoHandler struct {
	service *service.GrupoService
}

func NewGrupoHandler(service *service.GrupoService) *GrupoHandler {
	return &GrupoHandler{service: service}
}

func (h *GrupoHandler) Listar(c *gin.Context) {
	grupos, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToGrupoModels(grupos))
}

func (h *GrupoHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("grupoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	grupo, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToGrupoModel(grupo))
}

func (h *GrupoHandler) Adicionar(c *gin.Context) {
	var input dto.GrupoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	grupo := assembler.ToGrupoEntity(&input)
	if err := h.service.Save(grupo); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assembler.ToGrupoModel(grupo))
}

func (h *GrupoHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("grupoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.GrupoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	grupo, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	grupo.Nome = input.Nome
	if err := h.service.Save(grupo); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToGrupoModel(grupo))
}

func (h *GrupoHandler) Remover(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("grupoId"), 10, 64)
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

// ListarPermissoes lists permissions of a group
func (h *GrupoHandler) ListarPermissoes(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("grupoId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	grupo, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, assembler.ToPermissaoModels(grupo.Permissoes))
}

// AssociarPermissao adds a permission to a group
func (h *GrupoHandler) AssociarPermissao(c *gin.Context) {
	grupoID, _ := strconv.ParseUint(c.Param("grupoId"), 10, 64)
	permissaoID, _ := strconv.ParseUint(c.Param("permissaoId"), 10, 64)

	if err := h.service.AssociarPermissao(grupoID, permissaoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DesassociarPermissao removes a permission from a group
func (h *GrupoHandler) DesassociarPermissao(c *gin.Context) {
	grupoID, _ := strconv.ParseUint(c.Param("grupoId"), 10, 64)
	permissaoID, _ := strconv.ParseUint(c.Param("permissaoId"), 10, 64)

	if err := h.service.DesassociarPermissao(grupoID, permissaoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
