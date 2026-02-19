package repository

import (
	"time"

	"github.com/yurisasc/algafood-go/internal/domain/model"
	domainRepo "github.com/yurisasc/algafood-go/internal/domain/repository"
	"github.com/yurisasc/algafood-go/pkg/pagination"
	"gorm.io/gorm"
)

type pedidoRepositoryImpl struct {
	db *gorm.DB
}

// NewPedidoRepository creates a new PedidoRepository
func NewPedidoRepository(db *gorm.DB) *pedidoRepositoryImpl {
	return &pedidoRepositoryImpl{db: db}
}

func (r *pedidoRepositoryImpl) FindAll(filter *domainRepo.PedidoFilter, page *pagination.Pageable) (*pagination.Page[model.Pedido], error) {
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

	// Busca apenas dados básicos - relacionamentos serão populados via cache no serviço
	if err := query.
		Offset(page.Offset()).
		Limit(page.Size).
		Order("data_criacao DESC").
		Find(&pedidos).Error; err != nil {
		return nil, err
	}

	return pagination.NewPage(pedidos, total, page), nil
}

func (r *pedidoRepositoryImpl) FindByCodigo(codigo string) (*model.Pedido, error) {
	var pedido model.Pedido
	// Carrega apenas os itens - os outros relacionamentos serão populados via cache no serviço
	if err := r.db.
		Preload("Itens").
		Where("codigo = ?", codigo).
		First(&pedido).Error; err != nil {
		return nil, err
	}
	return &pedido, nil
}

func (r *pedidoRepositoryImpl) Save(pedido *model.Pedido) error {
	// Usa Omit para evitar que o GORM tente inserir/atualizar as associações
	// Apenas os IDs das foreign keys serão salvos
	return r.db.Omit("Restaurante", "Cliente", "FormaPagamento", "EnderecoEntrega.Cidade", "Itens.Produto").Save(pedido).Error
}

func (r *pedidoRepositoryImpl) IsPedidoGerenciadoPor(codigoPedido string, usuarioID uint64) (bool, error) {
	var count int64
	if err := r.db.Table("pedido p").
		Joins("JOIN restaurante_usuario_responsavel rur ON rur.restaurante_id = p.restaurante_id").
		Where("p.codigo = ? AND rur.usuario_id = ?", codigoPedido, usuarioID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
