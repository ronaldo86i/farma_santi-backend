package domain

import (
	"time"

	"github.com/google/uuid"
)

type Producto struct {
	Id                int
	NombreComercial   string
	FormaFarmaceutica string
	PrecioCompra      float64
	PrecioVenta       float64
	Estado            string
	Fotos             []string
	Stock             int64
	StockMin          int64
	CreatedAt         time.Time
	DeletedAt         *time.Time
}

type ProductoPrincipioActivoRequest struct {
	PrincipioActivoId int     `json:"principioActivoId"`
	Concentracion     float64 `json:"concentracion"`
	UnidadMedidaId    int     `json:"unidadMedidaId"`
}

type ProductoPrincipioActivo struct {
	Concentracion   float64             `json:"concentracion"`
	UnidadMedida    UnidadMedida        `json:"unidadMedida"`
	PrincipioActivo PrincipioActivoInfo `json:"principioActivo"`
}

type ProductRequest struct {
	NombreComercial      string                           `json:"nombreComercial"`
	PrincipiosActivos    []ProductoPrincipioActivoRequest `json:"principiosActivos"`
	FormaFarmaceuticaId  int                              `json:"formaFarmaceuticaId"`
	PrecioVenta          float64                          `json:"precioVenta"`
	StockMin             int64                            `json:"stockMin"`
	PresentacionId       int                              `json:"presentacionId"`
	UnidadesPresentacion int                              `json:"unidadesPresentacion"`
	Categorias           []int                            `json:"categorias"`
	LaboratorioId        int                              `json:"laboratorioId"`
}

type ProductoInfo struct {
	Id                   uuid.UUID    `json:"id"`
	NombreComercial      string       `json:"nombreComercial"`
	FormaFarmaceutica    string       `json:"formaFarmaceutica"`
	Laboratorio          string       `json:"laboratorio"`
	PrecioCompra         float64      `json:"precioCompra"`
	PrecioVenta          float64      `json:"precioVenta"`
	Stock                int64        `json:"stock"`
	StockMin             int64        `json:"stockMin"`
	Estado               string       `json:"estado"`
	Presentacion         Presentacion `json:"presentacion"`
	UnidadesPresentacion int          `json:"unidadesPresentacion"`
	UrlFoto              string       `json:"urlFoto,omitempty"`
	DeletedAt            *time.Time   `json:"deletedAt"`
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

type Presentacion struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type ProductoDetail struct {
	Id                   uuid.UUID                 `json:"id"`
	NombreComercial      string                    `json:"nombreComercial"`
	FormaFarmaceutica    FormaFarmaceutica         `json:"formaFarmaceutica"`
	Laboratorio          LaboratorioSimple         `json:"laboratorio"`
	Categorias           []Categoria               `json:"categorias"`
	PrincipiosActivos    []ProductoPrincipioActivo `json:"principiosActivos"`
	PrecioCompra         float64                   `json:"precioCompra"`
	PrecioVenta          float64                   `json:"precioVenta"`
	StockMin             int64                     `json:"stockMin"`
	Stock                int64                     `json:"stock"`
	Estado               string                    `json:"estado"`
	Presentacion         Presentacion              `json:"presentacion"`
	UnidadesPresentacion int                       `json:"unidadesPresentacion"`
	UrlFotos             []string                  `json:"urlFotos"`
	CreatedAt            time.Time                 `json:"createdAt"`
	DeletedAt            *time.Time                `json:"deletedAt"`
}

type ProductoSimple struct {
	Id                   uuid.UUID    `json:"id"`
	NombreComercial      string       `json:"nombreComercial"`
	Laboratorio          string       `json:"laboratorio,omitempty"`
	FormaFarmaceutica    string       `json:"formaFarmaceutica,omitempty"`
	Presentacion         Presentacion `json:"presentacion,omitempty"`
	UnidadesPresentacion int          `json:"unidadesPresentacion,omitempty"`
}

type ProductoId struct {
	Id string `json:"id"`
}

type VentaResponse struct {
	VentaId int64 `json:"ventaId"`
}
