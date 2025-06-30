package repository

import (
	"context"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrincipioActivoRepository struct {
	pool *pgxpool.Pool
}

func (p PrincipioActivoRepository) RegistrarPrincipioActivo(ctx context.Context, request *domain.PrincipioActivoRequest) error {
	query := `INSERT INTO principio_activo(nombre, descripcion) VALUES ($1, $2)`

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, query, request.Nombre, request.Descripcion)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Código 23505 = unique_violation
			return datatype.NewConflictError("El nombre del principio activo ya existe")
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (p PrincipioActivoRepository) ModificarPrincipioActivo(ctx context.Context, id *int, request *domain.PrincipioActivoRequest) error {
	query := `UPDATE principio_activo SET nombre = $1, descripcion = $2 WHERE id = $3`

	result, err := p.pool.Exec(ctx, query, request.Nombre, request.Descripcion, *id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Error de restricción UNIQUE violada
				return datatype.NewConflictError("El principio activo ya existe")
			}
		}
		return err
	}

	if result.RowsAffected() == 0 {
		return datatype.NewNotFoundError("No existe el principio activo")
	}
	return nil
}

func (p PrincipioActivoRepository) ListarPrincipioActivo(ctx context.Context) (*[]domain.PrincipioActivoInfo, error) {
	query := `SELECT id, nombre, descripcion FROM principio_activo ORDER BY nombre`
	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var lista = make([]domain.PrincipioActivoInfo, 0)
	for rows.Next() {
		var pa domain.PrincipioActivoInfo
		if err := rows.Scan(&pa.Id, &pa.Nombre, &pa.Descripcion); err != nil {
			return nil, err
		}
		lista = append(lista, pa)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &lista, nil
}

func (p PrincipioActivoRepository) ObtenerPrincipioActivoById(ctx context.Context, id *int) (*domain.PrincipioActivoDetail, error) {
	query := `SELECT id, nombre, descripcion FROM principio_activo WHERE id = $1`
	row := p.pool.QueryRow(ctx, query, *id)

	var detalle domain.PrincipioActivoDetail
	err := row.Scan(&detalle.Id, &detalle.Nombre, &detalle.Descripcion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, datatype.NewNotFoundError("El principio activo no existe")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &detalle, nil
}

func NewPrincipioActivoRepository(pool *pgxpool.Pool) *PrincipioActivoRepository {
	return &PrincipioActivoRepository{pool: pool}
}

var _ port.PrincipioActivoRepository = (*PrincipioActivoRepository)(nil)
