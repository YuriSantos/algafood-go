package repository

import (
	"time"

	"github.com/yurisasc/algafood-go/internal/domain/model"
	domainRepo "github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type PedidoRepositoryImpl struct {
	db *gorm.DB
}

// NewPedidoRepository creates a new PedidoRepository
func NewPedidoRepository(db *gorm.DB) *PedidoRepositoryImpl {
	return &PedidoRepositoryImpl{db: db}
}

func (r *PedidoRepositoryImpl) FindAll(filter *domainRepo.PedidoFilter, page *pagination.Pageable) (*pagination.Page[model.Pedido], error) {
	var pedidos []model.Pedido
	var total int64

	query := r.db.Model(&model.Pedido{})

	// Apply filters
	if filter != nil {
		if filter.ClienteID != nil {
			query = query.Where("usuario_cliente_id = ?", *filter.ClienteID)
		}
		if filter.RestauranteID != nil {
			query = query.Where("restaurante_id = ?", *filter.RestauranteID)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.DataCriacaoInicio != nil {
			t, _ := time.Parse("2006-01-02", *filter.DataCriacaoInicio)
			query = query.Where("data_criacao >= ?", t)
		}
		if filter.DataCriacaoFim != nil {
			t, _ := time.Parse("2006-01-02", *filter.DataCriacaoFim)
			query = query.Where("data_criacao <= ?", t)
		}
	}

	query.Count(&total)

	// Busca apenas dados básicos do pedido (sem Preload pesados)
	if err := query.
		Offset(page.Offset()).
		Limit(page.Size).
		Order("data_criacao DESC").
		Find(&pedidos).Error; err != nil {
		return nil, err
	}

	return pagination.NewPage(pedidos, total, page), nil
}

func (r *PedidoRepositoryImpl) FindByCodigo(codigo string) (*model.Pedido, error) {
	var pedido model.Pedido
	// Busca apenas os itens do pedido (necessário para cálculos)
	// Os outros relacionamentos serão populados via cache no serviço
	if err := r.db.
		Preload("Itens").
		Preload("Itens.Produto").
		Where("codigo = ?", codigo).
		First(&pedido).Error; err != nil {
		return nil, err
	}
	return &pedido, nil
}

// FindByCodigoCompleto busca o pedido com todos os relacionamentos (para casos específicos)
func (r *PedidoRepositoryImpl) FindByCodigoCompleto(codigo string) (*model.Pedido, error) {
	var pedido model.Pedido
	if err := r.db.
		Preload("Restaurante").
		Preload("Restaurante.Cozinha").
		Preload("Cliente").
		Preload("FormaPagamento").
		Preload("Itens").
		Preload("Itens.Produto").
		Preload("EnderecoEntrega.Cidade").
		Preload("EnderecoEntrega.Cidade.Estado").
		Where("codigo = ?", codigo).
		First(&pedido).Error; err != nil {
		return nil, err
	}
	return &pedido, nil
}

func (r *PedidoRepositoryImpl) Save(pedido *model.Pedido) error {
	return r.db.Save(pedido).Error
}

func (r *PedidoRepositoryImpl) IsPedidoGerenciadoPor(codigoPedido string, usuarioID uint64) (bool, error) {
	var count int64
	if err := r.db.Table("pedido p").
		Joins("JOIN restaurante_usuario_responsavel rur ON rur.restaurante_id = p.restaurante_id").
		Where("p.codigo = ? AND rur.usuario_id = ?", codigoPedido, usuarioID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
