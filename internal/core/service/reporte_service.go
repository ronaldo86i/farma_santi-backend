package service

import (
	"context"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"github.com/johnfercher/maroto/v2"
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
	"log"
	"time"
)

type ReporteService struct {
	usuarioRepository      port.UsuarioRepository
	clienteRepository      port.ClienteRepository
	loteProductoRepository port.LoteProductoRepository
	productoRepository     port.ProductoRepository
	compraRepository       port.CompraRepository
	ventaRepository        port.VentaRepository
}

func (r ReporteService) ReporteUsuariosPDF(ctx context.Context) (core.Document, error) {
	usuarios, err := r.usuarioRepository.ListarUsuarios(ctx)
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
			text.NewCol(12, "Reporte de usuarios", props.Text{
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

func (r ReporteService) ReporteClientesPDF(ctx context.Context) (core.Document, error) {
	clientes, err := r.clienteRepository.ObtenerListaClientes(ctx)
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
			text.NewCol(12, "Reporte de clientes", props.Text{
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

func (r ReporteService) ReporteComprasPDF(ctx context.Context) (core.Document, error) {
	compras, err := r.compraRepository.ObtenerListaCompras(ctx)
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
			text.NewCol(12, "Reporte de compras", props.Text{
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
		text.NewCol(2, "NIT de proveedor", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha y Hora", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Usuario", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "Total (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de lotes
	for _, c := range *compras {
		m.AddAutoRow(
			text.NewCol(2, c.Codigo.String, props.Text{Style: fontstyle.Normal, Right: 2, Bottom: 1, Align: align.Left}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%d", c.Proveedor.NIT), props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
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

func (r ReporteService) ReporteVentasPDF(ctx context.Context) (core.Document, error) {
	ventas, err := r.ventaRepository.ObtenerListaVentas(ctx)
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
		WithTitle("Reporte de ventas", true).
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
			text.NewCol(12, "Reporte de ventas", props.Text{
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
		text.NewCol(2, "CI/NIT", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Fecha y Hora", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Estado", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
		text.NewCol(2, "Cajero", props.Text{Style: fontstyle.Bold, Align: align.Center, Bottom: 2}).WithStyle(colStyle),
		text.NewCol(2, "Total (Bs)", props.Text{Style: fontstyle.Bold, Align: align.Center}).WithStyle(colStyle),
	)

	// Datos de ventas
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
		// Cuerpo de la tabla
		m.AddAutoRow(
			text.NewCol(2, c.Codigo.String, props.Text{Style: fontstyle.Normal, Left: 2, Bottom: 1, Align: align.Left}).WithStyle(colStyle),
			text.NewCol(2, fmt.Sprintf("%s", nitCi), props.Text{Style: fontstyle.Normal, Align: align.Left, Left: 2, BreakLineStrategy: breakline.EmptySpaceStrategy}).WithStyle(colStyle),
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

func (r ReporteService) ReporteInventarioPDF(ctx context.Context) (core.Document, error) {
	productos, err := r.productoRepository.ObtenerListaProductos(ctx)
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
			text.NewCol(12, "Reporte de inventario", props.Text{
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

func (r ReporteService) ReporteLotesProductosPDF(ctx context.Context) (core.Document, error) {
	lotes, err := r.loteProductoRepository.ObtenerListaLotesProductos(ctx)
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
			text.NewCol(12, "Reporte de lotes de productos", props.Text{
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
) *ReporteService {
	return &ReporteService{
		usuarioRepository:      usuarioRepository,
		clienteRepository:      clienteRepository,
		loteProductoRepository: loteProductoRepository,
		productoRepository:     productoRepository,
		compraRepository:       compraRepository,
		ventaRepository:        ventaRepository,
	}
}

var _ port.ReporteService = (*ReporteService)(nil)
