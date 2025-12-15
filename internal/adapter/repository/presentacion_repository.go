package repository

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PresentacionRepository struct {
	pool *pgxpool.Pool
}

func (p PresentacionRepository) ObtenerListaPresentaciones(ctx context.Context) (*[]domain.Presentacion, error) {
	query := `SELECT p.id,p.nombre FROM presentacion p ORDER BY p.nombre`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		log.Println("Error al consultar presentaciones:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	list := make([]domain.Presentacion, 0)
	for rows.Next() {
		var item domain.Presentacion
		err := rows.Scan(&item.Id, &item.Nombre)
		if err != nil {
			log.Println("Error al scannar lista presentaciones:", err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}
	return &list, nil
}

func NewPresentacionRepository(pool *pgxpool.Pool) *PresentacionRepository {
	return &PresentacionRepository{pool: pool}
}

var _ port.PresentacionRepository = (*PresentacionRepository)(nil)
