package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"time"
)

type LaboratorioRepository struct {
	db *database.DB
}

func (l LaboratorioRepository) ListarLaboratoriosDisponibles(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	query := "SELECT l.id, l.nombre,l.estado,l.direccion, l.created_at,l.deleted_at FROM laboratorio l WHERE l.estado = 'Activo' ORDER BY id"
	rows, err := l.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	var list []domain.LaboratorioInfo
	for rows.Next() {
		var lab domain.LaboratorioInfo
		if err := rows.Scan(&lab.Id, &lab.Nombre, &lab.Estado, &lab.Direccion, &lab.CreatedAt, &lab.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerError()
		}
		list = append(list, lab)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}
	if len(list) == 0 {
		return &[]domain.LaboratorioInfo{}, nil
	}
	return &list, nil
}

func (l LaboratorioRepository) ListarLaboratorios(ctx context.Context) (*[]domain.LaboratorioInfo, error) {
	query := "SELECT l.id, l.nombre,l.estado,l.direccion, l.created_at,l.deleted_at FROM laboratorio l ORDER BY id "
	rows, err := l.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	var list []domain.LaboratorioInfo
	for rows.Next() {
		var lab domain.LaboratorioInfo
		if err := rows.Scan(&lab.Id, &lab.Nombre, &lab.Estado, &lab.Direccion, &lab.CreatedAt, &lab.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerError()
		}
		list = append(list, lab)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}
	if len(list) == 0 {
		return &[]domain.LaboratorioInfo{}, nil
	}
	return &list, nil
}

func (l LaboratorioRepository) ObtenerLaboratorioById(ctx context.Context, id *int) (*domain.LaboratorioDetail, error) {
	query := "SELECT l.id, l.nombre, l.direccion, l.estado ,l.created_at, l.deleted_at FROM laboratorio l WHERE l.id = $1 ORDER BY l.id"
	row := l.db.Pool.QueryRow(ctx, query, id)

	var laboratorio domain.LaboratorioDetail
	if err := row.Scan(&laboratorio.Id, &laboratorio.Nombre, &laboratorio.Direccion, &laboratorio.Estado, &laboratorio.CreatedAt, &laboratorio.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Laboratorio no encontrado",
			}
		}
		return nil, datatype.NewInternalServerError()
	}

	return &laboratorio, nil
}

func (l LaboratorioRepository) RegistrarLaboratorio(ctx context.Context, laboratorioRequest *domain.LaboratorioRequest) error {
	// Verifica si el laboratorio existe y si fue eliminado o no
	queryCheck := `SELECT deleted_at FROM laboratorio WHERE nombre = $1 LIMIT 1`

	var deletedAt *time.Time
	err := l.db.Pool.QueryRow(ctx, queryCheck, laboratorioRequest.Nombre).Scan(&deletedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Error al ejecutar la consulta
		return datatype.NewInternalServerError()
	}

	// Si se encontró un laboratorio con ese nombre y no está eliminado
	if err == nil && deletedAt == nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Ya existe el laboratorio",
		}
	}

	// Insertar el nuevo laboratorio
	queryInsert := `INSERT INTO laboratorio(nombre, direccion, deleted_at) VALUES ($1, $2, NULL)`
	_, err = l.db.Pool.Exec(ctx, queryInsert, laboratorioRequest.Nombre, laboratorioRequest.Direccion)

	if err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func (l LaboratorioRepository) ModificarLaboratorio(ctx context.Context, id *int, laboratorioRequest *domain.LaboratorioRequest) error {
	query := `UPDATE laboratorio l SET nombre=$1,direccion=$2 WHERE l.id = $3`
	tx, err := l.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	_, err = tx.Exec(ctx, query, laboratorioRequest.Nombre, laboratorioRequest.Direccion, *id)
	if err != nil {
		// Revisar si el error es de tipo PgError y si la restricción es por NIT duplicado
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Error de restricción UNIQUE violada
				return &datatype.ErrorResponse{
					Code:    http.StatusConflict,
					Message: "Ya existe el laboratorio",
				}
			} else if pgErr.Code == "23503" {
				// Error de clave foránea, si es que la 'id' del proveedor no existe
				return &datatype.ErrorResponse{
					Code:    http.StatusConflict,
					Message: "No existe el laboratorio con ese 'id'",
				}
			}
		}

		// Manejo de error interno
		return datatype.NewInternalServerError()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func (l LaboratorioRepository) HabilitarLaboratorio(ctx context.Context, id *int) error {
	query := `UPDATE laboratorio l SET deleted_at = NULL, estado= 'Activo' WHERE l.id = $1`
	tx, err := l.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		return datatype.NewInternalServerError()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func (l LaboratorioRepository) DeshabilitarLaboratorio(ctx context.Context, id *int) error {
	query := `UPDATE laboratorio l SET deleted_at = CURRENT_TIMESTAMP, estado= 'Inactivo' WHERE l.id = $1`
	tx, err := l.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		return datatype.NewInternalServerError()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func NewLaboratorioRepository(db *database.DB) *LaboratorioRepository {
	return &LaboratorioRepository{db: db}
}

var _ port.LaboratorioRepository = (*LaboratorioRepository)(nil)
