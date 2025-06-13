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
	"log"
	"net/http"
)

type LoteProductoRepository struct {
	db *database.DB
}

func (l LoteProductoRepository) ModificarLoteProducto(ctx context.Context, id *int, request *domain.LoteProductoRequest) error {

	tx, err := l.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	query := `UPDATE lote_producto SET nro_lote=$1,fecha_vencimiento=$2,producto_id=$3 WHERE id = $4`
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
	var query = `SELECT lp.id,lp.nro_lote,lp.stock,lp.fecha_vencimiento,
       			json_build_object(
       				'id', p.id,
       				'nombreComercial',p.nombre_comercial
       			) AS producto
				FROM lote_producto lp
				LEFT JOIN public.producto p on p.id = lp.producto_id ORDER BY lp.fecha_vencimiento,p.nombre_comercial
				`
	rows, err := l.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.LoteProductoInfo
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

	if len(list) == 0 {
		return &[]domain.LoteProductoInfo{}, nil
	}

	return &list, nil
}

func (l LoteProductoRepository) RegistrarLoteProducto(ctx context.Context, request *domain.LoteProductoRequest) error {
	tx, err := l.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	query := `INSERT INTO lote_producto(nro_lote, fecha_vencimiento, producto_id) VALUES ($1, $2, $3)`
	fechaVencimiento := request.FechaVencimiento.Format("2006-01-02")
	// Insertamos el nuevo lote del producto
	_, err = tx.Exec(ctx, query, request.Lote, fechaVencimiento, request.ProductoId.String())
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

func (l LoteProductoRepository) ObtenerLoteProductoById(ctx context.Context, id *int) (*domain.LoteProductoDetail, error) {
	var query = `SELECT lp.id,lp.nro_lote,lp.stock,lp.fecha_vencimiento,
       			json_build_object(
       				'id', p.id,
       				'nombreComercial',p.nombre_comercial,
       				'concentracion',(p.concentracion || ' ' || um.abreviatura)::TEXT,
       				'formaFarmaceutica',ff.nombre::TEXT,
       				'laboratorio',l.nombre::TEXT,
       				'precioVenta',p.precio_venta,
       				'stock',p.stock,
       				'stockMin',p.stock_min,
       				'estado',p.estado,
       				'deletedAt',p.deleted_at
       			) AS producto
				FROM lote_producto lp
				LEFT JOIN public.producto p on p.id = lp.producto_id 
				LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
                LEFT JOIN unidad_medida um ON um.id = p.unidad_medida_id
				WHERE lp.id = $1
				ORDER BY lp.fecha_vencimiento,p.nombre_comercial 
				LIMIT 1`

	var item domain.LoteProductoDetail
	err := l.db.Pool.QueryRow(ctx, query, id).Scan(&item.Id, &item.Lote, &item.Stock, &item.FechaVencimiento, &item.Producto)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Lote de producto no encontrada",
			}
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &item, nil
}

func NewLoteProductoRepository(db *database.DB) *LoteProductoRepository {
	return &LoteProductoRepository{db: db}
}

var _ port.LoteProductoRepository = (*LoteProductoRepository)(nil)
