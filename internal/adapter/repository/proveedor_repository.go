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
	"log"
)

type ProveedorRepository struct {
	pool *pgxpool.Pool
}

func (p ProveedorRepository) HabilitarProveedor(ctx context.Context, id *int) error {
	query := `UPDATE proveedor p SET deleted_at = NULL,estado='Activo' WHERE p.id = $1`
	tx, err := p.pool.Begin(ctx)
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

func (p ProveedorRepository) DeshabilitarProveedor(ctx context.Context, id *int) error {
	query := `UPDATE proveedor p SET deleted_at=CURRENT_TIMESTAMP,estado='Inactivo' WHERE p.id = $1`
	tx, err := p.pool.Begin(ctx)
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

func (p ProveedorRepository) RegistrarProveedor(ctx context.Context, request *domain.ProveedorRequest) error {
	// Primero, verificamos si el proveedor ya existe en la base de datos
	queryCheck := `SELECT 1 FROM proveedor p WHERE p.nit = $1 LIMIT 1`

	// Ejecutar la consulta
	res, err := p.pool.Exec(ctx, queryCheck, request.NIT)
	if err != nil {
		log.Println("proveedor.registrar", err)
		// Si hay error en la consulta, retornamos un error de servicio no disponible
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	// Verificar si la consulta encontró alguna fila
	rowsAffected := res.RowsAffected()
	if rowsAffected > 0 {
		return datatype.NewConflictError("Ya existe el proveedor")
	}

	// Iniciar la transacción
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		// Si no se puede iniciar la transacción, retornar error interno
		return datatype.NewInternalServerErrorGeneric()
	}

	// Aseguramos que el rollback se ejecute si algo falla
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Preparar la consulta de inserción
	queryInsert := `INSERT INTO proveedor(nit, razon_social, representante, direccion, telefono, email, celular) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(ctx, queryInsert, request.NIT, request.RazonSocial, request.Representante, request.Direccion, request.Telefono, request.Email, request.Celular)
	if err != nil {

		// Si ocurre un error en la inserción, retornamos un error interno
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmar la transacción si no hubo errores
	err = tx.Commit(ctx)
	if err != nil {
		// Si ocurre un error al hacer commit, realizamos rollback
		_ = tx.Rollback(ctx)
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (p ProveedorRepository) ObtenerProveedorById(ctx context.Context, id *int) (*domain.ProveedorDetail, error) {
	var proveedor domain.ProveedorDetail
	query := `SELECT p.id,p.nit,p.razon_social,p.representante,p.direccion,p.telefono,p.celular,p.email,p.created_at,p.deleted_at FROM proveedor p WHERE p.id = $1 LIMIT 1`
	err := p.pool.QueryRow(ctx, query, *id).Scan(&proveedor.Id, &proveedor.NIT, &proveedor.RazonSocial, &proveedor.Representante, &proveedor.Direccion, &proveedor.Telefono, &proveedor.Celular, &proveedor.Email, &proveedor.CreatedAt, &proveedor.DeletedAt)
	if err != nil {
		log.Print(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, datatype.NewNotFoundError("No existe el proveedor")
		}
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	return &proveedor, nil
}

func (p ProveedorRepository) ListarProveedores(ctx context.Context) (*[]domain.ProveedorInfo, error) {
	query := `SELECT p.id, p.estado,p.nit, p.razon_social, p.representante, p.direccion, p.created_at, p.deleted_at FROM proveedor p`
	rows, err := p.pool.Query(ctx, query)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var proveedores = make([]domain.ProveedorInfo, 0)
	for rows.Next() {
		var proveedor domain.ProveedorInfo
		err = rows.Scan(&proveedor.Id, &proveedor.Estado, &proveedor.NIT, &proveedor.RazonSocial, &proveedor.Representante, &proveedor.Direccion, &proveedor.CreatedAt, &proveedor.DeletedAt)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		proveedores = append(proveedores, proveedor)
	}
	return &proveedores, nil
}

func (p ProveedorRepository) ModificarProveedor(ctx context.Context, id *int, request *domain.ProveedorRequest) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		// Si no se puede iniciar la transacción, retornar error interno
		return datatype.NewInternalServerErrorGeneric()
	}

	// Aseguramos que el rollback se ejecute si algo falla
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `UPDATE proveedor p SET nit=$1, razon_social=$2, representante=$3, direccion=$4, telefono=$5, celular=$6, email=$7 WHERE p.id = $8`
	_, err = tx.Exec(ctx, query, request.NIT, request.RazonSocial, request.Representante, request.Direccion, request.Telefono, request.Celular, request.Email, *id)

	if err != nil {
		// Revisar si el error es de tipo PgError y si la restricción es por NIT duplicado
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Error de restricción UNIQUE violada
				return datatype.NewConflictError("Ya existe un proveedor con ese NIT")
			} else if pgErr.Code == "23503" {
				// Error de clave foránea, si es que la 'id' del proveedor no existe
				return datatype.NewNotFoundError("No existe el proveedor con ese ID")
			}
		}

		// Manejo de error interno
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmar la transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	return nil
}

func NewProveedorRepository(pool *pgxpool.Pool) *ProveedorRepository {
	return &ProveedorRepository{pool: pool}
}

var _ port.ProveedorRepository = (*ProveedorRepository)(nil)
