package service

import (
	"context"

	"github.com/yurisasc/algafood-go/internal/domain/event"
	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/internal/infrastructure/eventbridge"
)

type FluxoPedidoService struct {
	pedidoRepo     repository.PedidoRepository
	pedidoSvc      *PedidoService
	eventPublisher eventbridge.EventPublisher
}

func NewFluxoPedidoService(
	pedidoRepo repository.PedidoRepository,
	pedidoSvc *PedidoService,
	eventPublisher eventbridge.EventPublisher,
) *FluxoPedidoService {
	return &FluxoPedidoService{
		pedidoRepo:     pedidoRepo,
		pedidoSvc:      pedidoSvc,
		eventPublisher: eventPublisher,
	}
}

func (s *FluxoPedidoService) Confirmar(codigoPedido string) error {
	pedido, err := s.pedidoSvc.FindByCodigo(codigoPedido)
	if err != nil {
		return err
	}

	if err := pedido.Confirmar(); err != nil {
		return exception.NewNegocioException(err.Error())
	}

	if err := s.pedidoRepo.Save(pedido); err != nil {
		return err
	}

	// Publish domain event
	evt := event.NewPedidoConfirmadoEvent(
		pedido.Codigo,
		pedido.Cliente.ID,
		pedido.Cliente.Nome,
		pedido.Cliente.Email,
		pedido.Restaurante.ID,
		pedido.Restaurante.Nome,
		pedido.ValorTotal,
		*pedido.DataConfirmacao,
	)

	if err := s.eventPublisher.Publish(context.Background(), evt); err != nil {
		// Log error but don't fail the operation
		// In production, consider using a retry mechanism or dead letter queue
		return nil
	}

	return nil
}

func (s *FluxoPedidoService) Cancelar(codigoPedido string) error {
	pedido, err := s.pedidoSvc.FindByCodigo(codigoPedido)
	if err != nil {
		return err
	}

	if err := pedido.Cancelar(); err != nil {
		return exception.NewNegocioException(err.Error())
	}

	if err := s.pedidoRepo.Save(pedido); err != nil {
		return err
	}

	// Publish domain event
	evt := event.NewPedidoCanceladoEvent(
		pedido.Codigo,
		pedido.Cliente.ID,
		pedido.Cliente.Nome,
		pedido.Cliente.Email,
		pedido.Restaurante.ID,
		pedido.Restaurante.Nome,
		pedido.ValorTotal,
		*pedido.DataCancelamento,
	)

	if err := s.eventPublisher.Publish(context.Background(), evt); err != nil {
		// Log error but don't fail the operation
		return nil
	}

	return nil
}

func (s *FluxoPedidoService) Entregar(codigoPedido string) error {
	pedido, err := s.pedidoSvc.FindByCodigo(codigoPedido)
	if err != nil {
		return err
	}

	if err := pedido.Entregar(); err != nil {
		return exception.NewNegocioException(err.Error())
	}

	if err := s.pedidoRepo.Save(pedido); err != nil {
		return err
	}

	// Publish domain event
	evt := event.NewPedidoEntregueEvent(
		pedido.Codigo,
		pedido.Cliente.ID,
		pedido.Cliente.Nome,
		pedido.Cliente.Email,
		pedido.Restaurante.ID,
		pedido.Restaurante.Nome,
		pedido.ValorTotal,
		*pedido.DataEntrega,
	)

	if err := s.eventPublisher.Publish(context.Background(), evt); err != nil {
		// Log error but don't fail the operation
		return nil
	}

	return nil
}
