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

type RolRepository struct {
	db *database.DB
}

func (r RolRepository) ModificarRol(ctx context.Context, id *int, rolRequestUpdate *domain.RolRequestUpdate) error {
	if rolRequestUpdate.DeletedAt != nil {
		*rolRequestUpdate.DeletedAt = time.Now()
	}
	query := `UPDATE negocio.rol SET nombre = $1, deleted_at = $2 WHERE id = $3`

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al iniciar la transacción: " + err.Error(),
		}
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)

	}(tx, ctx)

	_, err = tx.Exec(ctx, query, rolRequestUpdate.Nombre, rolRequestUpdate.DeletedAt, id)
	if err != nil {
		// Si el nombre ya existe (violación de restricción única)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &datatype.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Ya existe un rol con ese nombre",
			}
		}
		return datatype.NewInternalServerError()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func (r RolRepository) RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error {
	// Primero, verificamos si el rol ya existe en la base de datos
	queryCheck := `SELECT deleted_at FROM negocio.rol WHERE nombre = $1`
	var deletedAt *time.Time
	err := r.db.Pool.QueryRow(ctx, queryCheck, rolRequest.Nombre).Scan(&deletedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Error al ejecutar la consulta
		return datatype.NewInternalServerError()
	}

	// Si el rol existe y no está eliminado
	if err == nil && deletedAt == nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Ya existe el rol",
		}
	}

	// Si el rol existe pero está eliminado, lo restauramos
	if err == nil && deletedAt != nil {
		queryRestaurar := `
		UPDATE negocio.rol
		SET deleted_at = NULL
		WHERE nombre = $1;
		`
		tx, err := r.db.Pool.Begin(ctx)
		if err != nil {
			return datatype.NewInternalServerError()
		}
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx) // Aseguramos rollback si algo falla
		}(tx, ctx)

		// Restauramos el rol
		_, err = tx.Exec(ctx, queryRestaurar, rolRequest.Nombre)
		if err != nil {
			return datatype.NewInternalServerError()
		}

		// Confirmamos la transacción
		if err := tx.Commit(ctx); err != nil {
			return datatype.NewInternalServerError()
		}

		return nil
	}

	// Si el rol no existe, lo insertamos como un nuevo rol
	queryInsert := `
		INSERT INTO negocio.rol(nombre, deleted_at) 
		VALUES ($1, NULL) 
		RETURNING id, nombre, created_at, deleted_at;
	`

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al iniciar transacción: " + err.Error(),
		}
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			// No hacemos nada si Rollback falla
		}
	}(tx, ctx)

	rol := domain.Rol{
		Nombre:    rolRequest.Nombre,
		DeletedAt: nil,
	}

	err = tx.QueryRow(ctx, queryInsert, rol.Nombre).Scan(&rol.Id, &rol.Nombre, &rol.CreatedAt, &rol.DeletedAt)
	if err != nil {
		// Manejo de error por nombre duplicado (conflicto)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &datatype.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Ya existe un rol con ese nombre",
			}
		}
		return datatype.NewInternalServerError()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func (r RolRepository) ModificarEstadoRol(ctx context.Context, id *int) error {
	query := `
		UPDATE negocio.rol r 
		SET deleted_at = CASE 
			WHEN deleted_at IS NULL THEN now()
		    WHEN deleted_at IS NOT NULL THEN NULL
			END
		WHERE r.id = $1
	`
	tx, err := r.db.Pool.Begin(ctx)
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

func (r RolRepository) ListarRoles(ctx context.Context) (*[]domain.Rol, error) {
	query := "SELECT r.id, r.nombre, r.created_at, r.deleted_at FROM negocio.rol r ORDER BY created_at DESC "
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	var roles []domain.Rol
	for rows.Next() {
		var rol domain.Rol
		if err := rows.Scan(&rol.Id, &rol.Nombre, &rol.CreatedAt, &rol.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerError()
		}
		roles = append(roles, rol)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}
	if len(roles) == 0 {
		return &[]domain.Rol{}, nil
	}
	return &roles, nil
}

func (r RolRepository) ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error) {
	query := "SELECT r.id, r.nombre, r.created_at, r.deleted_at FROM negocio.rol r WHERE r.id = $1 ORDER BY r.id"
	row := r.db.Pool.QueryRow(ctx, query, id)

	var rol domain.Rol
	if err := row.Scan(&rol.Id, &rol.Nombre, &rol.CreatedAt, &rol.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Rol no encontrado",
			}
		}
		return nil, datatype.NewInternalServerError()
	}

	return &rol, nil
}

func NewRolRepository(db *database.DB) *RolRepository {
	return &RolRepository{db}
}

var _ port.RolRepository = (*RolRepository)(nil)
