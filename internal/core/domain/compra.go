package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DetalleCompraRequest struct {
	Cantidad       uint    `json:"cantidad"`
	PrecioCompra   float64 `json:"precioCompra"`
	PrecioVenta    float64 `json:"precioVenta"`
	LoteProductoId uint    `json:"loteProductoId"`
}

type CompraRequest struct {
	Comentario    string                 `json:"comentario,omitempty"`
	LaboratorioId uint                   `json:"laboratorioId"`
	UsuarioId     uint                   `json:"-"`
	Detalles      []DetalleCompraRequest `json:"detalles"`
}

type CompraInfo struct {
	Id          uint              `json:"id"`
	Codigo      pgtype.Text       `json:"codigo"`
	Comentario  string            `json:"comentario"`
	Estado      string            `json:"estado"`
	Total       float64           `json:"total"`
	Fecha       time.Time         `json:"fecha"`
	Laboratorio LaboratorioSimple `json:"laboratorio"`
	Usuario     UsuarioSimple     `json:"usuario"`
}

type DetalleCompraDAO struct {
	Id             uint      `json:"id"`
	Cantidad       uint      `json:"cantidad"`
	PrecioCompra   float64   `json:"precioCompra"`
	PrecioVenta    float64   `json:"precioVenta"`
	LoteProductoId uint      `json:"loteProductoId"`
	ProductoId     uuid.UUID `json:"productoId"`
}

type CompraDAO struct {
	Id            uint               `json:"id"`
	Comentario    string             `json:"comentario"`
	Estado        string             `json:"estado"`
	Total         float64            `json:"total"`
	LaboratorioId int                `json:"laboratorioId"`
	UsuarioId     uint               `json:"usuarioId"`
	Detalles      []DetalleCompraDAO `json:"detalles"`
}

type DetalleCompraDetail struct {
	Id           uint             `json:"id"`
	Cantidad     uint             `json:"cantidad"`
	PrecioCompra float64          `json:"precioCompra"`
	PrecioVenta  float64          `json:"precioVenta"`
	LoteProducto LoteProductoInfo `json:"loteProducto"`
}

type CompraDetail struct {
	Id          uint                  `json:"id"`
	Codigo      pgtype.Text           `json:"codigo"`
	Comentario  string                `json:"comentario"`
	Estado      string                `json:"estado"`
	Total       float64               `json:"total"`
	Fecha       time.Time             `json:"fecha"`
	DeletedAt   *time.Time            `json:"deletedAt"`
	Laboratorio LaboratorioSimple     `json:"laboratorio"`
	Usuario     UsuarioSimple         `json:"usuario"`
	Detalles    []DetalleCompraDetail `json:"detalles"`
}

type CompraId struct {
	Id uint `json:"id"`
}
