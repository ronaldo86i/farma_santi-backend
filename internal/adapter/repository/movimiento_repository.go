package repository

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MovimientoRepository struct {
	pool *pgxpool.Pool
}

func (m MovimientoRepository) ObtenerMovimientosKardex(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoKardex, error) {
	// Consulta Base
	baseQuery := `
		SELECT
			vk.id_fila,
			vk.fecha_movimiento,
			vk.tipo_movimiento,
			vk.documento,
			vk.codigo_lote,
			vk.fecha_vencimiento,
			vk.usuario,
			vk.cantidad_entrada,
			vk.cantidad_salida,
			vk.costo_unitario,
			vk.total_moneda
		FROM view_kardex vk 
	`

	var whereClauses []string
	var args []interface{}
	i := 1 // Contador para los placeholders de Postgres ($1, $2, etc.)

	// Construcción de Filtros (WHERE)
	if tipoMovimientoStr := filtros["tipoMovimiento"]; tipoMovimientoStr != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vk.tipo_movimiento = $%d", i))
		args = append(args, tipoMovimientoStr)
		i++
	}

	if productoIdStr := filtros["productoId"]; productoIdStr != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("vk.producto_id = $%d", i))
		args = append(args, productoIdStr)
		i++
	}

	// Ensamblar Query con WHERE
	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// ORDER BY (Fecha ascendente)
	query += " ORDER BY vk.fecha_movimiento ASC, vk.id_fila ASC"

	// LIMIT y OFFSET
	if limitStr := filtros["limit"]; limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return nil, datatype.NewBadRequestError("El valor de limit debe ser un número entero positivo")
		}
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, limit)
		i++

		if offsetStr := filtros["offset"]; offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				return nil, datatype.NewBadRequestError("El valor de offset debe ser un número entero no negativo")
			}
			query += fmt.Sprintf(" OFFSET $%d", i)
			args = append(args, offset)
			i++
		}
	}

	rows, err := m.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var movimientos = make([]domain.MovimientoKardex, 0)
	for rows.Next() {
		var mo domain.MovimientoKardex
		// Asegúrate de que el orden del Scan coincida exactamente con el SELECT
		err = rows.Scan(
			&mo.IdFila,
			&mo.FechaMovimiento,
			&mo.TipoMovimiento,
			&mo.Documento,
			&mo.CodigoLote,
			&mo.FechaVencimiento,
			&mo.Usuario,
			&mo.CantidadEntrada,
			&mo.CantidadSalida,
			&mo.CostoUnitario,
			&mo.TotalMoneda,
		)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		movimientos = append(movimientos, mo)
	}
	return &movimientos, nil
}

func (m MovimientoRepository) ObtenerListaMovimientos(ctx context.Context, filtros map[string]string) (*[]domain.MovimientoInfo, error) {

	baseQuery := `SELECT m.id, m.codigo, m.tipo, m.estado, m.fecha, m.usuario, m.total FROM view_movimiento_info m`

	var filters []string
	var args []interface{}
	i := 1

	if val, ok := filtros["fechaInicio"]; ok && val != "" {
		filters = append(filters, fmt.Sprintf("m.fecha >= $%d", i))
		args = append(args, val)
		i++
	}

	if val, ok := filtros["fechaFin"]; ok && val != "" {
		filters = append(filters, fmt.Sprintf("m.fecha <= $%d", i))
		args = append(args, val)
		i++
	}

	if val, ok := filtros["tipo"]; ok && val != "" {
		filters = append(filters, fmt.Sprintf("m.tipo = $%d", i))
		args = append(args, val)
		i++
	}
	query := baseQuery
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	query += " ORDER BY m.fecha DESC;"

	rows, err := m.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var movimientos = make([]domain.MovimientoInfo, 0)
	for rows.Next() {
		var mo domain.MovimientoInfo
		err = rows.Scan(&mo.Id, &mo.Codigo, &mo.Tipo, &mo.Estado, &mo.Fecha, &mo.Usuario, &mo.Total)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		movimientos = append(movimientos, mo)
	}
	return &movimientos, nil
}

func NewMovimientoRepository(pool *pgxpool.Pool) *MovimientoRepository {
	return &MovimientoRepository{pool: pool}
}

var _ port.MovimientoRepository = (*MovimientoRepository)(nil)
