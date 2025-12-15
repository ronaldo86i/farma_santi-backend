package service

import (
	"context"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/breakline"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type ReporteService struct {
	usuarioRepository      port.UsuarioRepository
	clienteRepository      port.ClienteRepository
	loteProductoRepository port.LoteProductoRepository
	productoRepository     port.ProductoRepository
	compraRepository       port.CompraRepository
	ventaRepository        port.VentaRepository
	movimientoRepository   port.MovimientoRepository
}

func (r ReporteService) ReporteComprasDetallePDF(ctx context.Context, compraId *int) (core.Document, error) {
	// 1. Validar Usuario
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	// 2. Validar ID de Compra
	if compraId == nil {
		return nil, datatype.NewBadRequestError("El ID de la compra es requerido")
	}

	// 3. Obtener Datos de la Compra (Cabecera y Detalles)
	compra, err := r.compraRepository.ObtenerCompraById(ctx, compraId)
	if err != nil {
		log.Println("Error obteniendo compra:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	if compra == nil {
		return nil, datatype.NewBadRequestError("Compra no encontrada")
	}

	// 4. Configurar PDF
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("FarmaSanti System", true).
		WithTitle(fmt.Sprintf("Nota_Compra_%d", *compraId), true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Vertical). // Vertical para notas de detalle
		Build()

	m := maroto.New(cfg)

	// --- HEADER DEL REPORTE ---
	err = m.RegisterHeader(
		// Fila 1: Logo y Título
		row.New(25).Add(
			image.NewFromFileCol(2, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85,
			}),
			text.NewCol(7, "NOTA DE COMPRA / INGRESO", props.Text{
				Top:    8,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   14,
				Family: fontfamily.Helvetica,
			}),
			text.NewCol(3, fmt.Sprintf("Generado:\n%s", time.Now().Format("02/01/2006 15:04")), props.Text{
				Top:   2,
				Align: align.Right,
				Size:  8,
			}),
		),
		// Fila 2: Espacio
		row.New(5),
		// Fila 3: Datos de la Compra (Proveedor, Fecha, Código)
		row.New(20).Add(
			text.NewCol(6, fmt.Sprintf("Proveedor: %s\nCódigo Doc: %s", compra.Laboratorio.Nombre, util.Text.Coalesce(&compra.Codigo.String)), props.Text{
				Top:   0,
				Align: align.Left,
				Size:  10,
				Style: fontstyle.Bold,
			}),
			text.NewCol(6, fmt.Sprintf("Fecha Compra: %s\nRegistrado por: %s", compra.Fecha.Format("02/01/2006 15:04"), compra.Usuario.Username), props.Text{
				Top:   0,
				Align: align.Right,
				Size:  10,
			}),
		),
		// Fila 4: Estado
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Estado: %s", compra.Estado), props.Text{
				Style: fontstyle.Italic,
				Size:  10,
				Align: align.Left,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf header:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// --- FOOTER ---
	_ = m.RegisterFooter(
		row.New(15).Add(
			text.NewCol(6, fmt.Sprintf("Emitido por: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   8,
				Family: fontfamily.Arial,
				Top:    5,
			}),
			// Total general en el pie de página para resaltar
			text.NewCol(6, fmt.Sprintf("TOTAL COMPRA: %.2f Bs", compra.Total), props.Text{
				Align: align.Right,
				Size:  11,
				Style: fontstyle.Bold,
				Top:   2,
			}),
		),
	)

	// --- ESTILOS DE TABLA ---
	headerStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}

	colStyle := &props.Cell{
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 200, Green: 200, Blue: 200},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.1,
	}

	// --- TABLA DE DETALLES ---
	m.AddAutoRow(
		text.NewCol(1, "N°", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(4, "Producto", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Lote / Venc.", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(1, "Cant.", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Costo U. (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Subtotal (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
	)

	// Iterar sobre los detalles de la compra (Productos)
	// Ajustado al modelo: CompraDetail -> Detalles -> LoteProducto -> Producto
	for i, detalle := range compra.Detalles {
		// Formatear info del lote
		loteInfo := detalle.LoteProducto.Lote
		if !detalle.LoteProducto.FechaVencimiento.IsZero() {
			loteInfo = fmt.Sprintf("%s\n%s", detalle.LoteProducto.Lote, detalle.LoteProducto.FechaVencimiento.Format("02/01/06"))
		}

		// Calcular subtotal de línea
		subtotal := float64(detalle.Cantidad) * detalle.PrecioCompra

		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Size: 8, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(4, detalle.LoteProducto.Producto.NombreComercial, props.Text{Size: 8, Align: align.Left, Left: 2, BreakLineStrategy: breakline.DashStrategy}).WithStyle(colStyle),
			text.NewCol(2, loteInfo, props.Text{Size: 8, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", detalle.Cantidad), props.Text{Size: 9, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%.2f", detalle.PrecioCompra), props.Text{Size: 9, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%.2f", subtotal), props.Text{Size: 9, Align: align.Right, Right: 2, Style: fontstyle.Bold}).WithStyle(colStyle),
		)
	}
	m.AddAutoRow(
		text.NewCol(10, "TOTAL GENERAL:", props.Text{
			Style: fontstyle.Bold,
			Align: align.Right,
			Right: 2,
			Size:  9,
		}).WithStyle(colStyle),
		text.NewCol(2, fmt.Sprintf("%.2f Bs", compra.Total), props.Text{
			Style: fontstyle.Bold,
			Align: align.Right,
			Right: 2,
			Size:  9,
		}).WithStyle(colStyle),
	)
	// Generar documento
	document, err := m.Generate()
	if err != nil {
		log.Println("Error generando PDF detalle compra:", err.Error())
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteKardexProductoPDF(ctx context.Context, productoId *uuid.UUID) (core.Document, error) {
	// 1. Validación de usuario
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	// 2. Obtener datos del producto
	if productoId == nil {
		return nil, datatype.NewBadRequestError("El ID del producto es requerido")
	}
	producto, err := r.productoRepository.ObtenerProductoById(ctx, productoId)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// 3. Obtener movimientos del Kardex
	// Se asume que el repositorio acepta un mapa con el ID en string
	movimientos, err := r.movimientoRepository.ObtenerMovimientosKardex(ctx, map[string]string{"productoId": productoId.String()})
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	if movimientos == nil || len(*movimientos) == 0 {
		return nil, datatype.NewBadRequestError("El producto no tiene movimientos registrados")
	}

	// --- CÁLCULO DE STOCK VENCIDO ---
	stockVencido := 0
	hoy := time.Now()
	// Normalizamos 'hoy' al inicio del día para comparaciones consistentes (opcional, pero recomendado)
	hoy = time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, hoy.Location())

	for _, m := range *movimientos {
		// Si tiene fecha de vencimiento y ya pasó
		if !m.FechaVencimiento.IsZero() && m.FechaVencimiento.Before(hoy) {
			if m.TipoMovimiento == "ENTRADA" {
				stockVencido += m.CantidadEntrada
			} else {
				stockVencido -= m.CantidadSalida
			}
		}
	}
	// Aseguramos que no sea negativo por inconsistencias de datos antiguos
	if stockVencido < 0 {
		stockVencido = 0
	}

	// 4. Configuración del PDF
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("FarmaSanti System", true).
		WithTitle(fmt.Sprintf("Kardex_%s", producto.Id), true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal). // Horizontal para mejor visualización de columnas
		Build()

	m := maroto.New(cfg)

	// --- HEADER ---
	err = m.RegisterHeader(
		row.New(25).Add(
			image.NewFromFileCol(2, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85,
			}),
			text.NewCol(8, "KARDEX FÍSICO VALORADO", props.Text{
				Top:    8,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
			text.NewCol(2, fmt.Sprintf("Generado:\n%s", time.Now().Format("02/01/2006\n15:04")), props.Text{
				Top:   2,
				Align: align.Right,
				Size:  9,
			}),
		),
		// Fila de información del producto con el Stock Vencido agregado
		row.New(20).Add(
			text.NewCol(6, fmt.Sprintf("Producto: %s\nID: %s", producto.NombreComercial, producto.Id.String()), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
				Style: fontstyle.Bold,
			}),
			text.NewCol(6, fmt.Sprintf("Laboratorio: %s\nStock Actual: %d  |  Vencidos: %d", producto.Laboratorio.Nombre, producto.Stock, stockVencido), props.Text{
				Top:   2,
				Align: align.Right,
				Size:  10,
			}),
		),
	)
	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// --- FOOTER ---
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Emitido por: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)

	// --- ESTILOS ---
	headerStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}

	colStyle := &props.Cell{
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 200, Green: 200, Blue: 200},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.1,
	}

	// --- TABLA ---
	m.AddAutoRow(
		text.NewCol(2, "Fecha", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Documento", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Tipo", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(2, "Lote / Venc.", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(1, "Entrada", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(1, "Salida", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(1, "Saldo", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
		text.NewCol(1, "Usuario", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(headerStyle),
	)

	// Lógica de cálculo de saldo progresivo y llenado de filas
	saldoAcumulado := 0

	for _, mov := range *movimientos {
		entrada := mov.CantidadEntrada
		salida := mov.CantidadSalida

		if mov.TipoMovimiento == "ENTRADA" {
			saldoAcumulado += entrada
		} else {
			saldoAcumulado -= salida
		}

		strEntrada := "-"
		if entrada > 0 {
			strEntrada = fmt.Sprintf("%d", entrada)
		}

		strSalida := "-"
		if salida > 0 {
			strSalida = fmt.Sprintf("%d", salida)
		}

		// Combinamos Lote y Vencimiento para ahorrar espacio
		loteInfo := mov.CodigoLote
		if !mov.FechaVencimiento.IsZero() {
			loteInfo = fmt.Sprintf("%s\n%s", mov.CodigoLote, mov.FechaVencimiento.Format("02/01/06"))
		}

		m.AddAutoRow(
			text.NewCol(2, mov.FechaMovimiento.Format("02/01/2006 15:04"), props.Text{Size: 9, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(2, mov.Documento, props.Text{Size: 9, Align: align.Left, Left: 2}).WithStyle(colStyle),
			text.NewCol(2, mov.TipoMovimiento, props.Text{Size: 8, Align: align.Left, Left: 2}).WithStyle(colStyle),
			text.NewCol(2, loteInfo, props.Text{Size: 9, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(1, strEntrada, props.Text{Size: 9, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(1, strSalida, props.Text{Size: 9, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", saldoAcumulado), props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(1, mov.Usuario, props.Text{Size: 8, Align: align.Center}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteMovimientosPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}
	movs, err := r.movimientoRepository.ObtenerListaMovimientos(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	if len(*movs) == 0 {
		return nil, datatype.NewBadRequestError("Reporte sin movimientos")
	}
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte_movimientos", true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de movimientos", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)
	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)

	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddAutoRow(
		text.NewCol(1, "Nro", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "CÓDIGO", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Tipo", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Usuario", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Total (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)
	// Datos de movimientos
	for i, mov := range *movs {
		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, Bottom: 2}).WithStyle(colStyle),
			text.NewCol(2, mov.Codigo.String, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, mov.Tipo, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, mov.Estado, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(2, mov.Fecha.Format("02/01/2006 15:04:05"), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, mov.Usuario.Username, props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%.2f", mov.Total), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteUsuariosPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	usuarios, err := r.usuarioRepository.ListarUsuarios(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte_usuarios", true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de usuarios", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)
	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)

	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddAutoRow(
		text.NewCol(1, "Nro", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "Usuario", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "CI", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(3, "Nombre completo", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha de registro", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de usuarios
	for i, u := range *usuarios {
		personaCi := fmt.Sprintf("%d%s", u.Persona.Ci, util.Text.Coalesce(u.Persona.Complemento))
		personaNombreCompleto := fmt.Sprintf("%s %s %s", u.Persona.Nombres, u.Persona.ApellidoPaterno, u.Persona.ApellidoMaterno)

		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, Bottom: 2}).WithStyle(colStyle),
			text.NewCol(2, u.Username, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, personaCi, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(3, personaNombreCompleto, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, u.Estado, props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, u.CreatedAt.Format("02/01/2006 15:04:05"), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteClientesPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	clientes, err := r.clienteRepository.ObtenerListaClientes(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Construcción del reporte pdf
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte de clientes", true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de clientes", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)
	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddRows(
		row.New(9).Add(
			text.NewCol(1, "Nro", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(2, "CI/NIT", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(1, "Tipo de documento", props.Text{Style: fontstyle.Bold, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(3, "Razón social", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(1, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(2, "Fecha de registro", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		),
	)

	// Datos de clientes
	for i, c := range *clientes {
		var nitCi string
		if c.NitCi != nil {
			nitCi += fmt.Sprintf("%d", *c.NitCi)
			if c.Complemento.Valid {
				nitCi += c.Complemento.String
			}
		} else {
			nitCi = "Sin NIT/CI"
		}

		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, Bottom: 1}).WithStyle(colStyle),
			text.NewCol(2, nitCi, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, c.Tipo, props.Text{Style: fontstyle.Normal, Align: align.Center}).WithStyle(colStyle),
			text.NewCol(3, c.RazonSocial, props.Text{Style: fontstyle.Normal, Left: 2}).WithStyle(colStyle),
			text.NewCol(1, c.Estado, props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
			text.NewCol(2, c.CreatedAt.Format("02/01/2006 15:04:05"), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteComprasPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}
	compras, err := r.compraRepository.ObtenerListaCompras(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Construcción del reporte pdf
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte de compras", true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de compras", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)

	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddAutoRow(
		text.NewCol(2, "Código", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "Laboratorio", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha y Hora", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Usuario", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "Total (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de lotes
	for _, c := range *compras {
		m.AddAutoRow(
			text.NewCol(2, c.Codigo.String, props.Text{Style: fontstyle.Normal, Right: 2, Bottom: 1, Align: align.Left}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%s", c.Laboratorio.Nombre), props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, c.Fecha.Format("02/01/2006 15:04:05"), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, c.Estado, props.Text{Style: fontstyle.Normal, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, c.Usuario.Username, props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%.2f", c.Total), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteVentasPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	ventas, err := r.ventaRepository.ObtenerListaVentas(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	if len(*ventas) == 0 {
		return nil, datatype.NewNotFoundError("Reporte vacío")
	}
	// --- LÓGICA DE VISIBILIDAD DE COLUMNA ---
	mostrarColumnaEstado := false
	estadoFiltro := filtros["estado"]
	tituloReporte := "Reporte de ventas"

	// Si NO hay filtro o es "Todos", debemos mostrar la columna porque los estados varían
	if estadoFiltro == "" || estadoFiltro == "Todos" {
		mostrarColumnaEstado = true
	} else {
		// Si hay un filtro específico, lo mostramos en el título y ocultamos la columna repetitiva
		tituloReporte = fmt.Sprintf("Reporte de ventas (%s)", estadoFiltro)
	}

	// --- AJUSTE DINÁMICO DE TAMAÑOS (Sistema de 12 columnas) ---
	// Distribución Base (Con Columna Estado):
	// Código(1) + CI(2) + Fecha(2) + Estado(1) + TipoPago(2) + Cajero(2) + Total(2) = 12
	colAnchoFecha := 2

	if !mostrarColumnaEstado {
		// Si ocultamos la columna estado (1 espacio), se lo sumamos a Fecha
		colAnchoFecha = 3
	}

	// Construcción del reporte pdf
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle(tituloReporte, true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85,
			}),
			text.NewCol(9, tituloReporte, props.Text{ // Usamos la variable tituloReporte
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)

	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}

	// --- CONSTRUCCIÓN DINÁMICA DE CABECERAS ---
	var headerCols []core.Col
	headerCols = append(headerCols, text.NewCol(2, "Cód.", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle))
	headerCols = append(headerCols, text.NewCol(1, "CI/NIT", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle))
	headerCols = append(headerCols, text.NewCol(colAnchoFecha, "Fecha", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle))

	// Columna Estado (Condicional: Si es Todos/Vacio -> MOSTRAR)
	if mostrarColumnaEstado {
		headerCols = append(headerCols, text.NewCol(1, "Est.", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle))
	}

	headerCols = append(headerCols, text.NewCol(2, "Forma Pago", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle))
	headerCols = append(headerCols, text.NewCol(2, "Cajero", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle))
	headerCols = append(headerCols, text.NewCol(2, "Total (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle))

	// Agregar cabecera
	m.AddAutoRow(headerCols...)

	// --- CONSTRUCCIÓN DINÁMICA DE FILAS DE DATOS ---
	for _, c := range *ventas {
		var nitCi string
		if c.Cliente.NitCi != nil {
			nitCi += fmt.Sprintf("%d", *c.Cliente.NitCi)
			if c.Cliente.Complemento.Valid {
				nitCi += c.Cliente.Complemento.String
			}
		} else {
			nitCi = "Sin NIT/CI"
		}

		var rowCols []core.Col
		rowCols = append(rowCols, text.NewCol(2, c.Codigo.String, props.Text{Style: fontstyle.Normal, Size: 8, Left: 2, Bottom: 1, Align: align.Left}).WithStyle(colStyle))
		rowCols = append(rowCols, text.NewCol(1, fmt.Sprintf("%s", nitCi), props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))
		rowCols = append(rowCols, text.NewCol(colAnchoFecha, c.Fecha.Format("02/01/06 15:04"), props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))

		if mostrarColumnaEstado {
			rowCols = append(rowCols, text.NewCol(1, c.Estado, props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))
		}

		rowCols = append(rowCols, text.NewCol(2, c.TipoPago, props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))
		rowCols = append(rowCols, text.NewCol(2, c.Usuario.Username, props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))
		rowCols = append(rowCols, text.NewCol(2, fmt.Sprintf("%.2f", c.Total), props.Text{Style: fontstyle.Normal, Size: 8, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle))

		m.AddAutoRow(rowCols...)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}
func (r ReporteService) ReporteInventarioPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	productos, err := r.productoRepository.ObtenerListaProductos(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Construcción del reporte pdf
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte de inventario", true).
		WithPageNumber(pageNumber).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(10).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de inventario", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)
	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddAutoRow(
		text.NewCol(1, "Nro", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(3, "Nombre comercial", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Laboratorio", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Forma", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
		text.NewCol(1, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(1, "Cantidad Mínima", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(1, "Cantidad", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(1, "Precio de compra", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Precio de Venta", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de ventas
	for i, p := range *productos {

		// Cuerpo de la tabla
		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Style: fontstyle.Normal, Right: 2, Bottom: 1, Align: align.Right}).WithStyle(colStyle),
			text.NewCol(3, p.NombreComercial, props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, p.Laboratorio, props.Text{Style: fontstyle.Normal, Align: align.Left, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, p.FormaFarmaceutica, props.Text{Style: fontstyle.Normal, Align: align.Center, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, p.Estado, props.Text{Style: fontstyle.Normal, Align: align.Center, Bottom: 1, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", p.StockMin), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", p.Stock), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%.2f", p.PrecioCompra), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%.2f", p.PrecioVenta), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func (r ReporteService) ReporteLotesProductosPDF(ctx context.Context, filtros map[string]string) (core.Document, error) {
	userId, ok := ctx.Value(util.ContextUserIdKey).(int)
	if !ok {
		return nil, datatype.NewStatusUnauthorizedError("Usuario no autorizado")
	}
	usuario, err := r.usuarioRepository.ObtenerUsuarioDetalle(ctx, &userId)
	if err != nil {
		return nil, err
	}

	lotes, err := r.loteProductoRepository.ObtenerListaLotesProductos(ctx, filtros)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Construcción del reporte pdf
	pageNumber := props.PageNumber{
		Pattern: "Página {current} de {total}",
		Place:   props.RightBottom,
		Family:  fontfamily.Arial,
		Style:   fontstyle.Normal,
		Size:    9,
	}

	cfg := config.NewBuilder().
		WithCreator("Maroto v2", true).
		WithTitle("Reporte de lotes de productos", true).
		WithPageNumber(pageNumber).
		WithTopMargin(20).
		WithLeftMargin(20).
		WithRightMargin(20).
		WithOrientation(orientation.Horizontal).
		Build()

	m := maroto.New(cfg)

	// Título
	err = m.RegisterHeader(
		row.New(20).Add(
			image.NewFromFileCol(1, "./public/Logo.png", props.Rect{
				Center:  true,
				Percent: 85, // Ajusta el tamaño de la imagen dentro de la celda (85% de altura)
			}),
			text.NewCol(9, "Reporte de lotes de productos", props.Text{
				Top:    5,
				Style:  fontstyle.Bold,
				Align:  align.Center,
				Size:   16,
				Family: fontfamily.Helvetica,
			}),
		),
		row.New(10).Add(
			text.NewCol(12, fmt.Sprintf("Fecha y Hora: %s", time.Now().Format("02/01/2006 15:04:05")), props.Text{
				Top:   2,
				Align: align.Left,
				Size:  10,
			}),
		),
	)

	if err != nil {
		log.Println("Error al construir pdf:", err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Footer con usuario
	_ = m.RegisterFooter(
		row.New(10).Add(
			text.NewCol(6, fmt.Sprintf("Usuario: %s", usuario.Username), props.Text{
				Align:  align.Left,
				Size:   9,
				Family: fontfamily.Arial,
			}),
		),
	)
	// Estilo de columna
	colStyle := &props.Cell{
		BackgroundColor: &props.Color{Red: 255, Green: 255, Blue: 255},
		BorderType:      border.Full,
		BorderColor:     &props.Color{Red: 0, Green: 0, Blue: 0},
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.2,
	}
	// Encabezado de tabla
	m.AddAutoRow(
		text.NewCol(1, "Nro", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(3, "Nombre de producto", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Lote", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha de vencimiento", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Laboratorio", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(1, "Cantidad (Unidades)", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(1, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de lotes
	for i, l := range *lotes {
		m.AddAutoRow(
			text.NewCol(1, fmt.Sprintf("%d", i+1), props.Text{Style: fontstyle.Normal, Right: 2, Bottom: 1, Align: align.Right, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(3, l.Producto.NombreComercial, props.Text{Style: fontstyle.Normal, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, l.Lote, props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, l.FechaVencimiento.Format("02/01/2006"), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(2, l.Producto.Laboratorio, props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, fmt.Sprintf("%d", l.Stock), props.Text{Style: fontstyle.Normal, Align: align.Right, Right: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
			text.NewCol(1, l.Estado, props.Text{Style: fontstyle.Normal, Right: 2, Align: align.Right, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
		)
	}

	document, err := m.Generate()
	if err != nil {
		return nil, datatype.NewInternalServerError("Error al generar archivo .pdf")
	}
	return document, nil
}

func NewReporteService(
	usuarioRepository port.UsuarioRepository,
	clienteRepository port.ClienteRepository,
	loteProductoRepository port.LoteProductoRepository,
	productoRepository port.ProductoRepository,
	compraRepository port.CompraRepository,
	ventaRepository port.VentaRepository,
	movimientoRepository port.MovimientoRepository,
) *ReporteService {
	return &ReporteService{
		usuarioRepository:      usuarioRepository,
		clienteRepository:      clienteRepository,
		loteProductoRepository: loteProductoRepository,
		productoRepository:     productoRepository,
		compraRepository:       compraRepository,
		ventaRepository:        ventaRepository,
		movimientoRepository:   movimientoRepository,
	}
}

var _ port.ReporteService = (*ReporteService)(nil)
