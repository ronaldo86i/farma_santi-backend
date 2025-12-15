package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RolRepository struct {
	pool *pgxpool.Pool
}

func (r RolRepository) HabilitarRol(ctx context.Context, id *int) error {
	query := `UPDATE rol r SET deleted_at = NULL, estado= 'Activo' WHERE r.id = $1`
	tx, err := r.pool.Begin(ctx)
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

func (r RolRepository) DeshabilitarRol(ctx context.Context, id *int) error {
	query := `
		UPDATE rol r 
		SET deleted_at = CURRENT_TIMESTAMP, estado= 'Inactivo'
		WHERE r.id = $1 AND r.nombre != 'ADMIN'
	`
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	ct, err := tx.Exec(ctx, query, *id)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		return datatype.NewConflictError("Conflicto al actualizar rol")
	}
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (r RolRepository) ModificarRol(ctx context.Context, id *int, rolRequestUpdate *domain.RolRequest) error {
	var cantidad int
	query := `SELECT count(*) AS cantidad FROM usuario_rol WHERE rol_id=$1`
	err := r.pool.QueryRow(ctx, query).Scan(&cantidad)
	if err != nil {
		log.Println("Error al consultar cantidad de usuarios:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	if cantidad > 0 {
		return datatype.NewBadRequestError("Rol no permitido para actualizar")
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query = `UPDATE rol SET nombre = $1 WHERE id = $2`
	ct, err := tx.Exec(ctx, query, rolRequestUpdate.Nombre, id)
	if err != nil {
		// Si el nombre ya existe (violación de restricción única)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("Ya existe el rol")
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("No existe el rol")
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (r RolRepository) RegistrarRol(ctx context.Context, rolRequest *domain.RolRequest) error {
	// Primero, verificamos si el rol ya existe en la base de datos
	query := `SELECT 1 FROM rol WHERE nombre = $1 LIMIT 1;`

	ct, err := r.pool.Exec(ctx, query, rolRequest.Nombre)

	if ct.RowsAffected() != 0 {
		return datatype.NewNotFoundError("Ya existe el rol")
	}

	var rolId uint
	query = `INSERT INTO rol(nombre, deleted_at) VALUES ($1, NULL) RETURNING id;`
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		log.Print("Error al iniciar transacción:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	err = tx.QueryRow(ctx, query, rolRequest.Nombre).Scan(&rolId)
	if err != nil {
		// Manejo de error por nombre duplicado (conflicto)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("Ya existe el rol")
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (r RolRepository) ListarRoles(ctx context.Context) (*[]domain.Rol, error) {
	query := "SELECT r.id, r.nombre,r.estado, r.created_at, r.deleted_at FROM rol r ORDER BY created_at DESC "
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var roles = make([]domain.Rol, 0)
	for rows.Next() {
		var rol domain.Rol
		if err := rows.Scan(&rol.Id, &rol.Nombre, &rol.Estado, &rol.CreatedAt, &rol.DeletedAt); err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		roles = append(roles, rol)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &roles, nil
}

func (r RolRepository) ObtenerRolById(ctx context.Context, id *int) (*domain.Rol, error) {
	query := "SELECT r.id, r.nombre, r.created_at, r.deleted_at FROM rol r WHERE r.id = $1 ORDER BY r.id"
	row := r.pool.QueryRow(ctx, query, id)

	var rol domain.Rol
	if err := row.Scan(&rol.Id, &rol.Nombre, &rol.CreatedAt, &rol.DeletedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Rol no encontrado")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &rol, nil
}

func NewRolRepository(pool *pgxpool.Pool) *RolRepository {
	return &RolRepository{pool: pool}
}

var _ port.RolRepository = (*RolRepository)(nil)
