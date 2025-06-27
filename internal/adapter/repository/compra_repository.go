package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"strings"
)

type CompraRepository struct {
	db *database.DB
}

func (c CompraRepository) ObtenerCompraById(ctx context.Context, id *int) (*domain.CompraDetail, error) {
	query := `SELECT v.id,v.estado,v.total,v.comentario,v.proveedor,v.usuario,v.created_at,deleted_at,v.detalles FROM view_compra_con_detalles v WHERE id = $1`
	var compra domain.CompraDetail
	err := c.db.Pool.QueryRow(ctx, query, *id).Scan(&compra.Id, &compra.Estado, &compra.Total, &compra.Comentario, &compra.Proveedor, &compra.Usuario, &compra.CreatedAt, &compra.DeletedAt, &compra.Detalles)
	if err != nil {
		log.Println("Error al obtener compra:", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Registro de compra no encontrada")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &compra, nil
}

func (c CompraRepository) RegistrarOrdenCompra(ctx context.Context, request *domain.CompraRequest) error {
	tx, err := c.db.Pool.Begin(ctx)
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
	var total float64
	for _, detalle := range request.Detalles {
		total += detalle.Precio * float64(detalle.Cantidad)
	}
	query := `INSERT INTO compra(usuario_id,proveedor_id,comentario,total) VALUES($1,$2,$3,$4) RETURNING id`

	var compraId uint
	err = tx.QueryRow(ctx, query, request.UsuarioId, request.ProveedorId, request.Comentario, total).Scan(&compraId)
	if err != nil {
		log.Println("Ha ocurrido un error al insertar la compra:", err.Error())
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	for _, detalle := range request.Detalles {
		query = `INSERT INTO detalle_compra(cantidad, precio, compra_id, lote_producto_id) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, query, detalle.Cantidad, detalle.Precio, compraId, detalle.LoteProductoId)
		if err != nil {
			log.Println("Ha ocurrido un error al insertar detalles de la compra:", err.Error())
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}

	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c CompraRepository) ModificarOrdenCompra(ctx context.Context, id *int, request *domain.CompraRequest) error {
	// Iniciar transacción
	tx, err := c.db.Pool.Begin(ctx)
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
	// Verificar si existe la compra
	var existe bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM compra WHERE id = $1 AND estado='Pendiente')`
	err = tx.QueryRow(ctx, checkQuery, *id).Scan(&existe)
	if err != nil {
		log.Println("Error al verificar existencia de compra:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	if !existe {
		return datatype.NewNotFoundError("La compra no existe o fue eliminada")
	}

	var total float64
	for _, detalle := range request.Detalles {
		total += detalle.Precio * float64(detalle.Cantidad)
	}
	// Consulta
	query := `UPDATE compra c SET comentario=$1, total=$2,proveedor_id=$3,updated_at=CURRENT_TIMESTAMP WHERE id=$4`
	_, err = tx.Exec(ctx, query, request.Comentario, total, request.ProveedorId, *id)
	if err != nil {
		log.Println("Ha ocurrido un error al modificar la compra:", err.Error())
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	query = `DELETE FROM detalle_compra WHERE compra_id = $1`
	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		log.Println("Ha ocurrido un error al modificar detalles de compra:", err.Error())
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	for _, detalle := range request.Detalles {
		query = `INSERT INTO detalle_compra(cantidad, precio, compra_id, lote_producto_id) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(ctx, query, detalle.Cantidad, detalle.Precio, *id, detalle.LoteProductoId)
		if err != nil {
			log.Println("Ha ocurrido un error al insertar detalles de la compra:", err.Error())
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}
	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c CompraRepository) AnularOrdenCompra(ctx context.Context, id *int) error {
	// Iniciar transacción
	tx, err := c.db.Pool.Begin(ctx)
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

	// Verificar si existe la compra y está en estado 'Pendiente'
	var existe bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM compra WHERE id = $1 AND estado = 'Pendiente')`
	err = tx.QueryRow(ctx, checkQuery, *id).Scan(&existe)
	if err != nil {
		log.Println("Error al verificar existencia de compra:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	if !existe {
		return datatype.NewNotFoundError("La compra no fue encontrada o no está en estado pendiente")
	}

	// Actualizar el estado a 'Anulado'
	updateQuery := `UPDATE compra SET estado = 'Anulado',deleted_at=CURRENT_TIMESTAMP WHERE id = $1`
	_, err = tx.Exec(ctx, updateQuery, *id)
	if err != nil {
		log.Println("Error al anular compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c CompraRepository) RegistrarCompra(ctx context.Context, id *int) error {
	var compra domain.CompraDAO

	query := `SELECT c.id, c.estado, c.total, c.comentario, c.proveedor_id, c.usuario_id, c.detalles 
	          FROM view_compras_detalle c 
	          WHERE id = $1 LIMIT 1`

	err := c.db.Pool.QueryRow(ctx, query, *id).Scan(&compra.Id, &compra.Estado, &compra.Total, &compra.Comentario, &compra.ProveedorId, &compra.UsuarioId, &compra.Detalles)
	if err != nil {
		log.Println("Error al consultar la compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	estado := strings.ToLower(compra.Estado)
	switch estado {
	case "completado":
		log.Println("La compra ya está completada, no se puede volver a registrar.")
		return datatype.NewConflictError("La compra ya fue registrada y completada")
	case "anulado":
		log.Println("La compra ya está anulada, no se puede registrar.")
		return datatype.NewConflictError("La compra ya fue anulada")
	}

	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		log.Println("Error al iniciar transacción:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Lock de compra
	lockCompraQuery := `SELECT id FROM compra WHERE id = $1 FOR UPDATE`
	_, err = tx.Exec(ctx, lockCompraQuery, *id)
	if err != nil {
		log.Println("Error al bloquear compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	productosId := make(map[string]bool)

	for _, detalle := range compra.Detalles {
		// Lock de lote_producto
		lockLoteQuery := `SELECT id FROM lote_producto WHERE id = $1 FOR UPDATE`
		_, err := tx.Exec(ctx, lockLoteQuery, detalle.LoteProductoId)
		if err != nil {
			log.Printf("Error al bloquear lote_producto %d: %v", detalle.LoteProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Lock de producto
		lockProductoQuery := `SELECT precio_compra, stock FROM producto WHERE id = $1 FOR UPDATE`
		var precioActual float64
		var stockActual uint

		err = tx.QueryRow(ctx, lockProductoQuery, detalle.ProductoId).Scan(&precioActual, &stockActual)
		if err != nil {
			log.Printf("Error al obtener precio_compra y stock del producto %d: %v", detalle.ProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		productosId[detalle.ProductoId.String()] = true

		// Actualizar stock en lote_producto
		updateLoteQuery := `UPDATE lote_producto SET stock = stock + $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateLoteQuery, detalle.Cantidad, detalle.LoteProductoId)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23514" && pgErr.ConstraintName == "check_fecha_vencimiento" {
					log.Printf("Violación de restricción check_fecha_vencimiento para lote %d", detalle.LoteProductoId)
					return datatype.NewBadRequestError("No se puede actualizar el stock: el lote ya está vencido.")
				}
			}
			log.Printf("Error al actualizar stock del lote %d: %v", detalle.LoteProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Actualizar stock en producto
		updateProductoQuery := `UPDATE producto SET stock = stock + $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateProductoQuery, detalle.Cantidad, detalle.ProductoId)
		if err != nil {
			log.Printf("Error al actualizar producto %d: %v", detalle.ProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}
	}

	// ✅ Actualizar estado antes de calcular promedio para incluir la misma compra
	updateEstadoQuery := `UPDATE compra SET estado = 'Completado', fecha = CURRENT_TIMESTAMP WHERE id = $1`
	_, err = tx.Exec(ctx, updateEstadoQuery, *id)
	if err != nil {
		log.Println("Error al actualizar estado de la compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	// 🔁 Calcular y actualizar precio promedio
	for productoIdStr := range productosId {
		productoUUID, err := uuid.Parse(productoIdStr)
		if err != nil {
			log.Printf("Error al parsear UUID del producto: %s", productoIdStr)
			return datatype.NewBadRequestError("ID de producto inválido")
		}

		query = `
			SELECT 
				COALESCE(
					SUM(dc.precio * dc.cantidad)::NUMERIC / NULLIF(SUM(dc.cantidad), 0),
					0
				) AS precio_promedio_ponderado
			FROM detalle_compra dc
			LEFT JOIN lote_producto lp ON dc.lote_producto_id = lp.id
			LEFT JOIN compra c ON dc.compra_id = c.id
			WHERE lp.producto_id = $1 AND c.estado = 'Completado'
		`

		var precioPromedio float64
		err = tx.QueryRow(ctx, query, productoUUID).Scan(&precioPromedio)
		if err != nil {
			log.Printf("Error al calcular precio promedio ponderado para producto %s: %v", productoIdStr, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		log.Printf("Producto %s nuevo precio promedio: %.2f", productoIdStr, precioPromedio)

		updateQuery := `UPDATE producto SET precio_compra = $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateQuery, precioPromedio, productoUUID)
		if err != nil {
			log.Printf("Error al actualizar precio_compra del producto %s: %v", productoIdStr, err)
			return datatype.NewInternalServerErrorGeneric()
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error al confirmar transacción:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	committed = true
	return nil
}

func (c CompraRepository) ObtenerListaCompras(ctx context.Context) (*[]domain.CompraInfo, error) {
	query := `SELECT c.id, c.comentario, c.estado, c.total, c.proveedor, c.usuario, c.created_at FROM view_listar_compras c`
	rows, err := c.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]domain.CompraInfo, 0)

	for rows.Next() {
		var item domain.CompraInfo
		err = rows.Scan(&item.Id, &item.Comentario, &item.Estado, &item.Total, &item.Proveedor, &item.Usuario, &item.CreatedAt)
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

func NewCompraRepository(db *database.DB) *CompraRepository {
	return &CompraRepository{db: db}
}

var _ port.CompraRepository = (*CompraRepository)(nil)
