package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/assembler"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/api/middleware"
	"github.com/yurisasc/algafood-go/internal/domain/model"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/internal/domain/service"
	"github.com/yurisasc/algafood-go/pkg/pagination"
)

type PedidoHandler struct {
	service      *service.PedidoService
	fluxoService *service.FluxoPedidoService
}

func NewPedidoHandler(service *service.PedidoService, fluxoService *service.FluxoPedidoService) *PedidoHandler {
	return &PedidoHandler{
		service:      service,
		fluxoService: fluxoService,
	}
}

func (h *PedidoHandler) Pesquisar(c *gin.Context) {
	page := pagination.NewPageableFromContext(c)

	// Build filter from query params
	filter := &repository.PedidoFilter{}

	if clienteIDStr := c.Query("clienteId"); clienteIDStr != "" {
		clienteID, _ := strconv.ParseUint(clienteIDStr, 10, 64)
		filter.ClienteID = &clienteID
	}
	if restauranteIDStr := c.Query("restauranteId"); restauranteIDStr != "" {
		restauranteID, _ := strconv.ParseUint(restauranteIDStr, 10, 64)
		filter.RestauranteID = &restauranteID
	}
	if dataInicio := c.Query("dataCriacaoInicio"); dataInicio != "" {
		filter.DataCriacaoInicio = &dataInicio
	}
	if dataFim := c.Query("dataCriacaoFim"); dataFim != "" {
		filter.DataCriacaoFim = &dataFim
	}
	if status := c.Query("status"); status != "" {
		statusPedido := model.StatusPedido(status)
		filter.Status = &statusPedido
	}

	result, err := h.service.Pesquisar(filter, page)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Convert to DTOs
	content := assembler.ToPedidoResumoModels(result.Content)
	response := pagination.Page[dto.PedidoResumoModel]{
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

func (h *PedidoHandler) Buscar(c *gin.Context) {
	codigoPedido := c.Param("codigoPedido")

	pedido, err := h.service.FindByCodigo(codigoPedido)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, assembler.ToPedidoModel(pedido))
}

func (h *PedidoHandler) Adicionar(c *gin.Context) {
	var input dto.PedidoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptionhandler.HandleValidationError(c, err)
		return
	}

	// Get authenticated user from context
	usuario, ok := middleware.GetCurrentUser(c)
	if !ok {
		exceptionhandler.HandleUnauthorized(c)
		return
	}

	pedido := assembler.ToPedidoEntity(&input, usuario.ID)
	if err := h.service.Emitir(pedido); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Reload to get full data
	pedido, _ = h.service.FindByCodigo(pedido.Codigo)
	c.JSON(http.StatusCreated, assembler.ToPedidoModel(pedido))
}

func (h *PedidoHandler) Confirmar(c *gin.Context) {
	codigoPedido := c.Param("codigoPedido")

	if err := h.fluxoService.Confirmar(codigoPedido); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PedidoHandler) Cancelar(c *gin.Context) {
	codigoPedido := c.Param("codigoPedido")

	if err := h.fluxoService.Cancelar(codigoPedido); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PedidoHandler) Entregar(c *gin.Context) {
	codigoPedido := c.Param("codigoPedido")

	if err := h.fluxoService.Entregar(codigoPedido); err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
