package repository

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovimientoRepository struct {
	pool *pgxpool.Pool
}

func (m MovimientoRepository) ObtenerListaMovimientos(ctx context.Context) (*[]domain.MovimientoInfo, error) {
	query := `SELECT  m.id,m.codigo,m.tipo,m.estado,m.fecha,m.usuario FROM view_movimiento_info m ORDER BY m.fecha DESC;`
	rows, err := m.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var movimientos = make([]domain.MovimientoInfo, 0)
	for rows.Next() {
		var mo domain.MovimientoInfo
		err = rows.Scan(&mo.Id, &mo.Codigo, &mo.Tipo, &mo.Estado, &mo.Fecha, &mo.Usuario)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		movimientos = append(movimientos, mo)
	}
	return &movimientos, nil
}

func NewMovimientoRepository(pool *pgxpool.Pool) *MovimientoRepository {
	return &MovimientoRepository{pool: pool}
}

var _ port.MovimientoRepository = (*MovimientoRepository)(nil)
