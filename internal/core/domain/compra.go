package domain

import (
	"github.com/google/uuid"
	"time"
)

type DetalleCompraRequest struct {
	Cantidad       uint    `json:"cantidad"`
	Precio         float64 `json:"precio"`
	LoteProductoId uint    `json:"loteProductoId"`
}

type CompraRequest struct {
	Comentario  string                 `json:"comentario,omitempty"`
	ProveedorId uint                   `json:"proveedorId"`
	UsuarioId   uint                   `json:"-"`
	Detalles    []DetalleCompraRequest `json:"detalles"`
}

type CompraInfo struct {
	Id         uint            `json:"id"`
	Codigo     string          `json:"codigo"`
	Comentario string          `json:"comentario"`
	Estado     string          `json:"estado"`
	Total      float64         `json:"total"`
	Fecha      time.Time       `json:"fecha"`
	Proveedor  ProveedorSimple `json:"proveedor"`
	Usuario    UsuarioSimple   `json:"usuario"`
}

type DetalleCompraDAO struct {
	Id             uint      `json:"id"`
	Cantidad       uint      `json:"cantidad"`
	Precio         float64   `json:"precio"`
	LoteProductoId uint      `json:"loteProductoId"`
	ProductoId     uuid.UUID `json:"productoId"`
}

type CompraDAO struct {
	Id          uint               `json:"id"`
	Comentario  string             `json:"comentario"`
	Estado      string             `json:"estado"`
	Total       float64            `json:"total"`
	ProveedorId int                `json:"proveedorId"`
	UsuarioId   uint               `json:"usuarioId"`
	Detalles    []DetalleCompraDAO `json:"detalles"`
}

type DetalleCompraDetail struct {
	Id           uint             `json:"id"`
	Cantidad     uint             `json:"cantidad"`
	Precio       float64          `json:"precio"`
	LoteProducto LoteProductoInfo `json:"loteProducto"`
}
type CompraDetail struct {
	Id         uint                  `json:"id"`
	Comentario string                `json:"comentario"`
	Estado     string                `json:"estado"`
	Total      float64               `json:"total"`
	Fecha      time.Time             `json:"fecha"`
	DeletedAt  *time.Time            `json:"deletedAt"`
	Proveedor  ProveedorSimple       `json:"proveedor"`
	Usuario    UsuarioSimple         `json:"usuario"`
	Detalles   []DetalleCompraDetail `json:"detalles"`
}
