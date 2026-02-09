package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yurisasc/algafood-go/internal/api/dto"
	"github.com/yurisasc/algafood-go/internal/api/exceptionhandler"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
)

type EstatisticaHandler struct {
	vendaQueryRepo repository.VendaQueryRepository
}

func NewEstatisticaHandler(vendaQueryRepo repository.VendaQueryRepository) *EstatisticaHandler {
	return &EstatisticaHandler{vendaQueryRepo: vendaQueryRepo}
}

func (h *EstatisticaHandler) ConsultarVendasDiarias(c *gin.Context) {
	filter := &repository.VendaDiariaFilter{}

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

	timeOffset := c.DefaultQuery("timeOffset", "+00:00")

	vendas, err := h.vendaQueryRepo.ConsultarVendasDiarias(filter, timeOffset)
	if err != nil {
		exceptionhandler.HandleError(c, err)
		return
	}

	// Convert to DTOs
	models := make([]dto.VendaDiariaModel, len(vendas))
	for i, v := range vendas {
		models[i] = dto.VendaDiariaModel{
			Data:          v.Data,
			TotalVendas:   v.TotalVendas,
			TotalFaturado: v.TotalFaturado,
		}
	}

	c.JSON(http.StatusOK, models)
}
