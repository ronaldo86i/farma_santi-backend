package domain

import (
	"github.com/google/uuid"
	"time"
)

type LoteProductoRequest struct {
	Lote             string    `json:"lote"`
	FechaVencimiento time.Time `json:"fechaVencimiento"`
	ProductoId       uuid.UUID `json:"productoId"`
}

type LoteProductoInfo struct {
	Id               int            `json:"id"`
	Lote             string         `json:"lote"`
	FechaVencimiento time.Time      `json:"fechaVencimiento"`
	Stock            int            `json:"stock"`
	Estado           string         `json:"estado"`
	Producto         ProductoSimple `json:"producto"`
}

type LoteProductoDetail struct {
	Id               int          `json:"id"`
	Lote             string       `json:"lote"`
	FechaVencimiento time.Time    `json:"fechaVencimiento"`
	Stock            int          `json:"stock"`
	Estado           string       `json:"estado"`
	Producto         ProductoInfo `json:"producto"`
}

type LoteProductoSimple struct {
	Id               int       `json:"id"`
	Lote             string    `json:"lote"`
	FechaVencimiento time.Time `json:"fechaVencimiento"`
}
