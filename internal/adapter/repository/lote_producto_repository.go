package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
)

type LoteProductoRepository struct {
	pool *pgxpool.Pool
}

func (l LoteProductoRepository) ActualizarLotesVencidos(ctx context.Context) error {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		log.Println("Error al iniciar transacción de lotes vencidos:", err)
		return datatype.NewInternalServerError("Error al iniciar transacción de lotes vencidos")
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `UPDATE lote_producto lp SET estado = 'Vencido' WHERE CURRENT_TIMESTAMP >= fecha_vencimiento AND estado != 'Vencido'`
	_, err = tx.Exec(ctx, query)
	if err != nil {
		log.Println("Error al actualizar lotes vencidos:", err)
		return datatype.NewInternalServerError("Error al actualizar lotes vencidos")
	}
	// Confirmamos la transacción
	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error al confirmar transacción de lotes vencidos:", err)
		return datatype.NewInternalServerError("Error al confirmar transacción de lotes vencidos")
	}
	committed = true
	return nil
}

func (l LoteProductoRepository) ListarLotesProductosByProductoId(ctx context.Context, productoId *uuid.UUID) (*[]domain.LoteProductoSimple, error) {
	query := `SELECT lp.id,lp.lote,lp.fecha_vencimiento FROM lote_producto lp WHERE lp.producto_id = $1`
	rows, err := l.pool.Query(ctx, query, productoId.String())
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al obtener la lista")
	}
	defer rows.Close()
	var list = make([]domain.LoteProductoSimple, 0)
	for rows.Next() {
		var item domain.LoteProductoSimple
		err := rows.Scan(&item.Id, &item.Lote, &item.FechaVencimiento)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}
	return &list, nil
}

func (l LoteProductoRepository) ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error {
	var tieneCompra int
	query := `SELECT 1 FROM detalle_compra dc WHERE dc.lote_producto_id=$1 LIMIT 1`
	err := l.pool.QueryRow(ctx, query, id).Scan(&tieneCompra)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println("Error en el servidor:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	if tieneCompra == 1 {
		return datatype.NewConflictError("El lote tiene compras asociadas, no es posible modificar")
	}
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	query = `UPDATE lote_producto SET lote=$1,fecha_vencimiento=$2,producto_id=$3 WHERE id = $4`
	fechaVencimiento := request.FechaVencimiento.Format("2006-01-02")
	// Insertamos el nuevo lote del producto
	_, err = tx.Exec(ctx, query, request.Lote, fechaVencimiento, request.ProductoId.String(), *id)
	if err != nil {
		// Detectar errores específicos de PostgreSQL
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				// unique_violation
				return datatype.NewConflictError("Ya existe un lote con esos datos para este producto")
			case "23503":
				// foreign_key_violation
				return datatype.NewBadRequestError("El producto seleccionado no existe")
			case "23514":
				// check_violation
				return datatype.NewBadRequestError("Fecha de vencimiento inválida o fuera de rango permitido")
			default:
				// Otro error de base de datos
				return datatype.NewInternalServerErrorGeneric()
			}
		}

		// Si no es error de Postgres conocido
		return datatype.NewInternalServerErrorGeneric()
	}
	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (l LoteProductoRepository) ListarLotesProductos(ctx context.Context) (*[]domain.LoteProductoInfo, error) {
	var query = `SELECT lp.id,lp.lote,lp.stock,lp.fecha_vencimiento,lp.producto FROM view_lotes_con_productos lp`
	rows, err := l.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al obtener la lista")
	}
	defer rows.Close()
	var list = make([]domain.LoteProductoInfo, 0)
	for rows.Next() {
		var item domain.LoteProductoInfo
		err := rows.Scan(&item.Id, &item.Lote, &item.Stock, &item.FechaVencimiento, &item.Producto)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &list, nil
}

func (l LoteProductoRepository) RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	query := `INSERT INTO lote_producto(lote, fecha_vencimiento, producto_id) VALUES ($1, $2, $3)`
	fechaVencimiento := request.FechaVencimiento.Format("2006-01-02")
	// Insertamos el nuevo lote del producto
	_, err = tx.Exec(ctx, query, request.Lote, fechaVencimiento, request.ProductoId.String())
	if err != nil {
		// Detectar errores específicos de PostgreSQL
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return datatype.NewConflictError("Ya existe un lote con esos datos para este producto")
			case "23503":
				return datatype.NewBadRequestError("El producto seleccionado no existe")
			case "P0001":
				// Código para excepciones lanzadas por RAISE EXCEPTION en trigger
				if strings.Contains(pgErr.Message, "fecha de vencimiento") {
					return datatype.NewBadRequestError("La fecha de vencimiento no puede ser menor que la fecha actual")
				}
				return datatype.NewInternalServerErrorGeneric()
			default:
				return datatype.NewInternalServerErrorGeneric()
			}
		}

		// Si no es error de Postgres conocido
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (l LoteProductoRepository) ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error) {
	var query = `SELECT lp.id,lp.lote,lp.stock,lp.fecha_vencimiento,lp.producto FROM obtener_lote_by_id($1) lp`

	var item domain.LoteProductoDetail
	err := l.pool.QueryRow(ctx, query, *id).Scan(&item.Id, &item.Lote, &item.Stock, &item.FechaVencimiento, &item.Producto)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Lote de producto no encontrada")
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &item, nil
}

func NewLoteProductoRepository(pool *pgxpool.Pool) *LoteProductoRepository {
	return &LoteProductoRepository{pool: pool}
}

var _ port.LoteProductoRepository = (*LoteProductoRepository)(nil)
