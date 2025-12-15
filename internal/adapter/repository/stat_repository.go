package repository

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatRepository struct {
	pool *pgxpool.Pool
}

func (r StatRepository) ObtenerEstadisticasDashboard(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// 1. Totales Generales (Ventas y Compras)
	// Usamos CTEs o subconsultas para ser eficientes
	queryTotals := `
		SELECT
			(SELECT COALESCE(SUM(total), 0) FROM venta WHERE estado = 'Realizada') as total_ventas,
			(SELECT COUNT(*) FROM venta WHERE estado = 'Realizada') as cant_ventas,
			(SELECT COALESCE(SUM(total), 0) FROM compra WHERE estado = 'Completado') as total_compras,
			(SELECT COUNT(*) FROM compra WHERE estado != 'Completado') as cant_compras
	`

	err := r.pool.QueryRow(ctx, queryTotals).Scan(
		&stats.TotalVentas,
		&stats.CantidadVentas,
		&stats.TotalCompras,
		&stats.CantidadCompras,
	)
	if err != nil {
		log.Println("Error obteniendo totales dashboard:", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// 2. Ventas Diarias (Últimos 7 días) para gráficas
	// Ajusta 'YYYY-MM-DD' según tu dialecto SQL si no es Postgres (Postgres usa TO_CHAR)
	queryDaily := `
		SELECT 
			TO_CHAR(fecha, 'YYYY-MM-DD') as dia, 
			COALESCE(SUM(total), 0) as total
		FROM venta
		WHERE estado = 'Realizada' 
		  AND fecha >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY dia
		ORDER BY dia
	`

	rows, err := r.pool.Query(ctx, queryDaily)
	if err != nil {
		log.Println("Error obteniendo ventas diarias:", err)
		// No retornamos error fatal, devolvemos stats parciales si falla esto
		stats.VentasDiarias = []domain.VentaDiaria{}
		return stats, nil
	}
	defer rows.Close()

	var diarias []domain.VentaDiaria
	for rows.Next() {
		var d domain.VentaDiaria
		if err := rows.Scan(&d.Fecha, &d.Total); err == nil {
			diarias = append(diarias, d)
		}
	}
	stats.VentasDiarias = diarias

	return stats, nil
}

func (s StatRepository) ObtenerTopProductosVendidos(ctx context.Context, _ map[string]string) (*[]domain.ProductoStat, error) {
	fullHostname := ctx.Value("fullHostname").(string)
	fullHostname = fmt.Sprintf("%s%s", fullHostname, "/uploads/productos")
	query := `
		SELECT 
			p.id, 
			p.nombre_comercial, 
			ARRAY(
				SELECT $1 || '/' || p.id || '/' || foto
				FROM unnest(p.fotos) AS foto
			) AS fotos,
			CAST(COALESCE(SUM(dv.cantidad), 0) AS INTEGER) as total_vendido
		FROM detalle_venta dv
		INNER JOIN venta v ON v.id = dv.venta_id
		INNER JOIN lote_producto lp ON dv.lote_id = lp.id
		INNER JOIN producto p ON p.id = lp.producto_id
		WHERE v.estado = 'Realizada'
		GROUP BY p.id, p.nombre_comercial, p.fotos
		ORDER BY total_vendido DESC
		LIMIT 10;
	`

	rows, err := s.pool.Query(ctx, query, fullHostname)
	if err != nil {
		log.Println("Error al obtener top productos:", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var lista = make([]domain.ProductoStat, 0)

	for rows.Next() {
		var item domain.ProductoStat
		err := rows.Scan(
			&item.Id,
			&item.NombreComercial,
			&item.Fotos,
			&item.Cantidad,
		)
		if err != nil {
			log.Println("Error escaneando top productos:", err)
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		lista = append(lista, item)
	}

	return &lista, nil
}

func NewStatRepository(pool *pgxpool.Pool) *StatRepository {
	return &StatRepository{pool: pool}
}

var _ port.StatRepository = (*StatRepository)(nil)
