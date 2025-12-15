package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LaboratorioRepository struct {
	pool *pgxpool.Pool
}

func (l LaboratorioRepository) ListarLaboratoriosDisponibles(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	query := "SELECT l.id, l.nombre,l.estado,l.direccion, l.created_at,l.deleted_at FROM laboratorio l WHERE l.estado = 'Activo' ORDER BY id"
	rows, err := l.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var list = make([]domain.LaboratorioInfo, 0)
	for rows.Next() {
		var lab domain.LaboratorioInfo
		if err := rows.Scan(&lab.Id, &lab.Nombre, &lab.Estado, &lab.Direccion, &lab.CreatedAt, &lab.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, lab)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &list, nil
}

func (l LaboratorioRepository) ListarLaboratorios(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	query := "SELECT l.id, l.nombre,l.estado,l.direccion, l.created_at,l.deleted_at FROM laboratorio l ORDER BY id "
	rows, err := l.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var list = make([]domain.LaboratorioInfo, 0)
	for rows.Next() {
		var lab domain.LaboratorioInfo
		if err := rows.Scan(&lab.Id, &lab.Nombre, &lab.Estado, &lab.Direccion, &lab.CreatedAt, &lab.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, lab)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &list, nil
}

func (l LaboratorioRepository) ObtenerLaboratorioById(ctx context.Context, id *int) (*domain.LaboratorioDetail, error) {
	query := "SELECT l.id, l.nombre, l.direccion, l.estado ,l.created_at, l.deleted_at, l.telefono, l.celular, l.email, l.representante FROM laboratorio l WHERE l.id = $1 ORDER BY l.id"
	row := l.pool.QueryRow(ctx, query, id)

	var lab domain.LaboratorioDetail
	if err := row.Scan(&lab.Id, &lab.Nombre, &lab.Direccion, &lab.Estado, &lab.CreatedAt, &lab.DeletedAt, &lab.Telefono, &lab.Celular, &lab.Email, &lab.Representante); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Laboratorio no encontrado")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &lab, nil
}

func (l LaboratorioRepository) RegistrarLaboratorio(ctx context.Context, request *domain.LaboratorioRequest) error {
	// Verifica si el laboratorio existe y si fue eliminado o no
	queryCheck := `SELECT deleted_at FROM laboratorio WHERE nombre = $1 LIMIT 1`

	var deletedAt *time.Time
	err := l.pool.QueryRow(ctx, queryCheck, request.Nombre).Scan(&deletedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Error al ejecutar la consulta
		return datatype.NewInternalServerErrorGeneric()
	}

	// Si se encontró un laboratorio con ese nombre y no está eliminado
	if err == nil && deletedAt == nil {
		return datatype.NewConflictError("Ya existe el laboratorio")
	}

	// Insertar el nuevo laboratorio
	queryInsert := `INSERT INTO laboratorio(nombre, direccion, deleted_at,telefono,email,celular,representante) VALUES ($1, $2, NULL, $3,$4,$5,$6)`
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, queryInsert, request.Nombre, request.Direccion, request.Telefono, request.Email, request.Celular, request.Representante)
	if err != nil {
		log.Println("Error al insertar Laboratorio:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (l LaboratorioRepository) ModificarLaboratorio(ctx context.Context, id *int, request *domain.LaboratorioRequest) error {
	query := `UPDATE laboratorio l SET nombre=$1,direccion=$2,email=$3,celular=$4,telefono=$5,representante=$6 WHERE l.id = $7`
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, request.Nombre, request.Direccion, request.Email, request.Celular, request.Telefono, request.Representante, *id)
	if err != nil {
		// Revisar si el error es de tipo PgError y si la restricción es por NIT duplicado
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Error de restricción UNIQUE violada
				return datatype.NewConflictError("Ya existe el laboratorio")
			} else if pgErr.Code == "23503" {
				// Error de clave foránea, si es que la 'id' del proveedor no existe
				return datatype.NewConflictError("No existe el laboratorio con ese 'id'")
			}
		}

		// Manejo de error interno
		return datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (l LaboratorioRepository) HabilitarLaboratorio(ctx context.Context, id *int) error {
	query := `UPDATE laboratorio l SET deleted_at = NULL, estado= 'Activo' WHERE l.id = $1`
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	return nil
}

func (l LaboratorioRepository) DeshabilitarLaboratorio(ctx context.Context, id *int) error {
	query := `UPDATE laboratorio l SET deleted_at = CURRENT_TIMESTAMP, estado= 'Inactivo' WHERE l.id = $1`
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func NewLaboratorioRepository(pool *pgxpool.Pool) *LaboratorioRepository {
	return &LaboratorioRepository{pool: pool}
}

var _ port.LaboratorioRepository = (*LaboratorioRepository)(nil)
