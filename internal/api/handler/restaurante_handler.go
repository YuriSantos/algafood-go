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

type RestauranteHandler struct {
	service *service.RestauranteService
}

func NewRestauranteHandler(service *service.RestauranteService) *RestauranteHandler {
	return &RestauranteHandler{service: service}
}

func (h *RestauranteHandler) Listar(c *gin.Context) {
	restaurantes, err := h.service.FindAll()
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToRestauranteResumoModels(restaurantes))
}

func (h *RestauranteHandler) Buscar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	restaurante, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToRestauranteModel(restaurante))
}

func (h *RestauranteHandler) Adicionar(c *gin.Context) {
	var input dto.RestauranteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	restaurante := assembler.ToRestauranteEntity(&input)
	if err := h.service.Save(restaurante); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Reload to get related entities
	restaurante, _ = h.service.FindByID(restaurante.ID)
	c.JSON(http.StatusCreated, assembler.ToRestauranteModel(restaurante))
}

func (h *RestauranteHandler) Atualizar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var input dto.RestauranteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	restaurante, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Update fields from input
	updated := assembler.ToRestauranteEntity(&input)
	restaurante.Nome = updated.Nome
	restaurante.TaxaFrete = updated.TaxaFrete
	restaurante.CozinhaID = updated.CozinhaID
	if input.Endereco != nil {
		restaurante.Endereco = updated.Endereco
	}

	if err := h.service.Save(restaurante); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Reload
	restaurante, _ = h.service.FindByID(restaurante.ID)
	c.JSON(http.StatusOK, assembler.ToRestauranteModel(restaurante))
}

func (h *RestauranteHandler) Ativar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err := h.service.Ativar(id); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) Inativar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err := h.service.Inativar(id); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) AtivarEmMassa(c *gin.Context) {
	var input dto.AtivacaoRestauranteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	if err := h.service.AtivarEmMassa(input.IDs); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) InativarEmMassa(c *gin.Context) {
	var input dto.AtivacaoRestauranteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	if err := h.service.InativarEmMassa(input.IDs); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) Abrir(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err := h.service.Abrir(id); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) Fechar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	if err := h.service.Fechar(id); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListarFormasPagamento lists payment methods of a restaurant
func (h *RestauranteHandler) ListarFormasPagamento(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	restaurante, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToFormaPagamentoModels(restaurante.FormasPagamento))
}

func (h *RestauranteHandler) AssociarFormaPagamento(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	formaPagamentoID, _ := strconv.ParseUint(c.Param("formaPagamentoId"), 10, 64)

	if err := h.service.AssociarFormaPagamento(restauranteID, formaPagamentoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) DesassociarFormaPagamento(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	formaPagamentoID, _ := strconv.ParseUint(c.Param("formaPagamentoId"), 10, 64)

	if err := h.service.DesassociarFormaPagamento(restauranteID, formaPagamentoID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListarResponsaveis lists responsible users of a restaurant
func (h *RestauranteHandler) ListarResponsaveis(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	restaurante, err := h.service.FindByID(id)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToUsuarioModels(restaurante.Responsaveis))
}

func (h *RestauranteHandler) AssociarResponsavel(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	usuarioID, _ := strconv.ParseUint(c.Param("usuarioId"), 10, 64)

	if err := h.service.AssociarResponsavel(restauranteID, usuarioID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *RestauranteHandler) DesassociarResponsavel(c *gin.Context) {
	restauranteID, _ := strconv.ParseUint(c.Param("restauranteId"), 10, 64)
	usuarioID, _ := strconv.ParseUint(c.Param("usuarioId"), 10, 64)

	if err := h.service.DesassociarResponsavel(restauranteID, usuarioID); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
