package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type ClienteRepository struct {
	pool *pgxpool.Pool
}

func (c ClienteRepository) ObtenerListaClientes(ctx context.Context) (*[]domain.ClienteInfo, error) {
	query := `SELECT c.id,c.nit_ci,c.complemento,c.tipo,c.razon_social,c.estado FROM cliente c`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]domain.ClienteInfo, 0)
	for rows.Next() {
		var item domain.ClienteInfo
		err = rows.Scan(&item.Id, &item.NitCi, &item.Complemento, &item.Tipo, &item.RazonSocial, &item.Estado)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &list, nil
}

func (c ClienteRepository) ObtenerClienteById(ctx context.Context, id *int) (*domain.ClienteDetail, error) {
	query := `SELECT c.id,c.nit_ci,c.complemento,c.tipo,c.razon_social,c.estado,c.email,c.telefono,c.created_at,c.created_at FROM cliente c WHERE c.id=$1`
	var cliente domain.ClienteDetail
	err := c.pool.QueryRow(ctx, query, *id).Scan(&cliente.Id, &cliente.NitCi, &cliente.Complemento, &cliente.Tipo, &cliente.RazonSocial, &cliente.Estado, &cliente.Email, &cliente.Telefono, &cliente.CreatedAt, &cliente.DeletedAt)
	if err != nil {
		log.Println("Error al obtener cliente:", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Registro de cliente no encontrado")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &cliente, nil
}

func (c ClienteRepository) RegistrarCliente(ctx context.Context, request *domain.ClienteRequest) (*int, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Normalizar nit_ci: si es 0, convertir a NULL
	var nitCi interface{}
	if request.NitCi != nil && *request.NitCi == 0 {
		nitCi = nil
	} else {
		nitCi = request.NitCi
	}

	var query string
	var params []interface{}

	switch request.Tipo {
	case "NIT":
		query = `INSERT INTO cliente(nit_ci, tipo, razon_social, email, estado, telefono) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		params = []interface{}{nitCi, request.Tipo, request.RazonSocial, request.Email, "Activo", request.Telefono}

	case "CI":
		query = `INSERT INTO cliente(nit_ci, complemento, tipo, razon_social, email, estado, telefono) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
		params = []interface{}{nitCi, request.Complemento, request.Tipo, request.RazonSocial, request.Email, "Activo", request.Telefono}
	default:
		log.Printf("Tipo de cliente desconocido: %s", request.Tipo)
		return nil, datatype.NewBadRequestError("Tipo de cliente no válido")
	}

	var clienteId int
	err = tx.QueryRow(ctx, query, params...).Scan(&clienteId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error PG al registrar cliente: %s (constraint: %s, code: %s)", pgErr.Message, pgErr.ConstraintName, pgErr.Code)

			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "idx_unique_ci_complemento" {
					return nil, datatype.NewConflictError("Ya existe un cliente con ese CI y complemento")
				} else if pgErr.ConstraintName == "idx_unique_nit" {
					return nil, datatype.NewConflictError("Ya existe un cliente con ese NIT")
				}
			case "23514": // check_violation
				return nil, datatype.NewBadRequestError(fmt.Sprintf("Violación de regla de negocio: %s", pgErr.ConstraintName))
			case "23502": // not_null_violation
				return nil, datatype.NewBadRequestError(fmt.Sprintf("Falta un campo obligatorio: %s", pgErr.ColumnName))
			case "22P02": // invalid_text_representation (por ejemplo, error al insertar texto donde se espera número)
				return nil, datatype.NewBadRequestError("Tipo de dato inválido")
			}
		}

		log.Printf("Error inesperado al registrar cliente: %v", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error al confirmar transacción de cliente: %v", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return &clienteId, nil
}

func (c ClienteRepository) ModificarClienteById(ctx context.Context, id *int, request *domain.ClienteRequest) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Obtener tipo de cliente
	var tipoCliente string
	query := `SELECT tipo FROM cliente WHERE id = $1`
	err = tx.QueryRow(ctx, query, *id).Scan(&tipoCliente)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Cliente con ID %d no encontrado", *id)
			return datatype.NewNotFoundError("Cliente no encontrado")
		}
		log.Printf("Error al consultar tipo de cliente con id %d: %v", *id, err)
		return datatype.NewInternalServerErrorGeneric()
	}

	// Normalizar nit_ci: si es 0, tratarlo como NULL
	var nitCi interface{}
	if request.NitCi != nil && *request.NitCi == 0 {
		nitCi = nil
	} else {
		nitCi = request.NitCi
	}

	var updateQuery string
	var params []interface{}

	switch tipoCliente {
	case "NIT":
		updateQuery = `UPDATE cliente SET nit_ci = $1, razon_social = $2, email = $3, telefono = $4 WHERE id = $5`
		params = []interface{}{nitCi, request.RazonSocial, request.Email, request.Telefono, *id}

	case "CI":
		updateQuery = `UPDATE cliente SET nit_ci = $1, complemento = $2, razon_social = $3, email = $4, telefono = $5 WHERE id = $6`
		params = []interface{}{nitCi, request.Complemento, request.RazonSocial, request.Email, request.Telefono, *id}

	default:
		log.Printf("Tipo de cliente desconocido: %s", tipoCliente)
		return datatype.NewBadRequestError("Tipo de cliente no válido")
	}

	_, err = tx.Exec(ctx, updateQuery, params...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("Error PG al registrar cliente: %s (constraint: %s, code: %s)", pgErr.Message, pgErr.ConstraintName, pgErr.Code)

			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "idx_unique_ci_complemento" {
					return datatype.NewConflictError("Ya existe un cliente con ese CI y complemento")
				} else if pgErr.ConstraintName == "idx_unique_nit" {
					return datatype.NewConflictError("Ya existe un cliente con ese NIT")
				}
			case "23514": // check_violation
				return datatype.NewBadRequestError(fmt.Sprintf("Violación de regla de negocio: %s", pgErr.ConstraintName))
			case "23502": // not_null_violation
				return datatype.NewBadRequestError(fmt.Sprintf("Falta un campo obligatorio: %s", pgErr.ColumnName))
			case "22P02": // invalid_text_representation (por ejemplo, error al insertar texto donde se espera número)
				return datatype.NewBadRequestError("Tipo de dato inválido")
			}
		}

		log.Printf("Error inesperado al registrar cliente: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error al confirmar transacción de cliente con id %d: %v", *id, err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c ClienteRepository) HabilitarCliente(ctx context.Context, id *int) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `UPDATE cliente SET estado='Activo',deleted_at=NULL WHERE id=$1`
	ct, err := tx.Exec(ctx, query, *id)
	if err != nil {
		log.Printf("Error al deshabilitar el cliente con id %d: %v", *id, err)
		return datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("No existe el cliente")
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error al confirmar transacción:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c ClienteRepository) DeshabilitarCliente(ctx context.Context, id *int) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `UPDATE cliente SET estado='Inactivo',deleted_at=CURRENT_TIMESTAMP WHERE id=$1`
	ct, err := tx.Exec(ctx, query, *id)
	if err != nil {
		log.Printf("Error al deshabilitar el cliente con id %d: %v", *id, err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("No existe el cliente")
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error al confirmar transacción:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func NewClienteRepository(pool *pgxpool.Pool) *ClienteRepository {
	return &ClienteRepository{pool: pool}
}

var _ port.ClienteRepository = (*ClienteRepository)(nil)
