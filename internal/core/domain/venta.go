package domain

import (
	"time"
)

type VentaRequest struct {
	UsuarioId uint                  `json:"-"`
	ClienteId uint                  `json:"clienteId"`
	Detalles  []DetalleVentaRequest `json:"detalles"`
}

type DetalleVentaRequest struct {
	ProductoId string `json:"productoId"`
	Cantidad   uint   `json:"cantidad"`
}

type VentaInfo struct {
	Id        uint          `json:"id"`
	Codigo    string        `json:"codigo"`
	Usuario   UsuarioSimple `json:"usuario"`
	Cliente   ClienteSimple `json:"cliente"`
	Fecha     time.Time     `json:"fecha"`
	Estado    string        `json:"estado"`
	DeletedAt *time.Time    `json:"deletedAt"`
}

type VentaDetail struct {
	VentaInfo
	Detalles []DetalleVentaDetail `json:"detalles"`
}

type DetalleVentaDetail struct {
	Id       uint           `json:"id"`
	Producto ProductoSimple `json:"producto"`
	Lotes    []VentaLote    `json:"lotes"`
	Cantidad uint           `json:"cantidad"`
	Precio   float64        `json:"precio"`
	Total    float64        `json:"total"`
}
type VentaLoteProducto struct {
	Id         int    `json:"id"`
	Cantidad   uint   `json:"cantidad"`
	ProductoId string `json:"productoId"`
}

type VentaLote struct {
	Id               int    `json:"id"`
	Lote             string `json:"lote"`
	Cantidad         int    `json:"cantidad"`
	FechaVencimiento string `json:"fechaVencimiento"`
}

type VentaLoteProductoDAO struct {
	Id          uint
	Stock       uint
	PrecioVenta float64
}
