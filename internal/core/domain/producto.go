package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id                int
	NombreComercial   string
	NombreGenerico    string
	Concentracion     int64
	FormaFarmaceutica string
	PrecioCompra      int64
	PrecioVenta       float64
	Estado            string
	Fotos             []string
	Stock             int64
	StockMin          int64
	CreatedAt         time.Time
	DeletedAt         *time.Time
}

type ProductRequest struct {
	NombreComercial     string  `json:"nombreComercial"`
	NombreGenerico      string  `json:"nombreGenerico"`
	Concentracion       int     `json:"concentracion"`
	UnidadMedidaId      int     `json:"unidadMedidaId"`
	FormaFarmaceuticaId int     `json:"formaFarmaceuticaId"`
	PrecioVenta         float64 `json:"precioVenta"`
	StockMin            int64   `json:"stockMin"`
	Categorias          []int   `json:"categorias"`
	LaboratorioId       int     `json:"laboratorioId"`
}

type ProductInfo struct {
	Id                uuid.UUID `json:"id"`
	NombreComercial   string    `json:"nombreComercial"`
	NombreGenerico    string    `json:"nombreGenerico"`
	Concentracion     string    `json:"concentracion"`
	FormaFarmaceutica string    `json:"formaFarmaceutica"`
	Laboratorio       string    `json:"laboratorio"`
	PrecioVenta       float64   `json:"precioVenta"`
	Stock             int64     `json:"stock"`
	StockMin          int64     `json:"stockMin"`
	Estado            string    `json:"estado"`
	UrlFoto           string    `json:"urlFoto"`
}

type UnidadMedida struct {
	Id          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Abreviatura string `json:"abreviatura,omitempty"`
}

type FormaFarmaceutica struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}
