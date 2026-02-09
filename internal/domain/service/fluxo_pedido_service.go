package service

import (
	"github.com/yurisasc/algafood-go/internal/domain/exception"
	"github.com/yurisasc/algafood-go/internal/domain/repository"
)

type FluxoPedidoService struct {
	pedidoRepo repository.PedidoRepository
	pedidoSvc  *PedidoService
}

func NewFluxoPedidoService(pedidoRepo repository.PedidoRepository, pedidoSvc *PedidoService) *FluxoPedidoService {
	return &FluxoPedidoService{
		pedidoRepo: pedidoRepo,
		pedidoSvc:  pedidoSvc,
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

	// TODO: Publish domain event PedidoConfirmadoEvent

	return s.pedidoRepo.Save(pedido)
}

func (s *FluxoPedidoService) Cancelar(codigoPedido string) error {
	pedido, err := s.pedidoSvc.FindByCodigo(codigoPedido)
	if err != nil {
		return err
	}

	if err := pedido.Cancelar(); err != nil {
		return exception.NewNegocioException(err.Error())
	}

	// TODO: Publish domain event PedidoCanceladoEvent

	return s.pedidoRepo.Save(pedido)
}

func (s *FluxoPedidoService) Entregar(codigoPedido string) error {
	pedido, err := s.pedidoSvc.FindByCodigo(codigoPedido)
	if err != nil {
		return err
	}

	if err := pedido.Entregar(); err != nil {
		return exception.NewNegocioException(err.Error())
	}

	return s.pedidoRepo.Save(pedido)
}
