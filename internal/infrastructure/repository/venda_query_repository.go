package repository

import (
	domainRepo "github.com/yurisasc/algafood-go/internal/domain/repository"
	"gorm.io/gorm"
)

type vendaQueryRepositoryImpl struct {
	db *gorm.DB
}

// NewVendaQueryRepository creates a new VendaQueryRepository
func NewVendaQueryRepository(db *gorm.DB) *vendaQueryRepositoryImpl {
	return &vendaQueryRepositoryImpl{db: db}
}

func (r *vendaQueryRepositoryImpl) ConsultarVendasDiarias(filter *domainRepo.VendaDiariaFilter, timeOffset string) ([]domainRepo.VendaDiaria, error) {
	var vendas []domainRepo.VendaDiaria

	query := `
		SELECT
			DATE(CONVERT_TZ(p.data_criacao, '+00:00', ?)) as data,
			COUNT(p.id) as total_vendas,
			SUM(p.valor_total) as total_faturado
		FROM pedido p
		WHERE p.status IN ('CONFIRMADO', 'ENTREGUE')
	`
	args := []interface{}{timeOffset}

	if filter != nil {
		if filter.RestauranteID != nil {
			query += " AND p.restaurante_id = ?"
			args = append(args, *filter.RestauranteID)
		}
		if filter.DataCriacaoInicio != nil {
			query += " AND p.data_criacao >= ?"
			args = append(args, *filter.DataCriacaoInicio)
		}
		if filter.DataCriacaoFim != nil {
			query += " AND p.data_criacao <= ?"
			args = append(args, *filter.DataCriacaoFim)
		}
	}

	query += " GROUP BY DATE(CONVERT_TZ(p.data_criacao, '+00:00', ?)) ORDER BY data"
	args = append(args, timeOffset)

	if err := r.db.Raw(query, args...).Scan(&vendas).Error; err != nil {
		return nil, err
	}

	return vendas, nil
}
