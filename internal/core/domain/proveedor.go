package domain

import "time"

type Proveedor struct {
	Id            int        `json:"id"`
	NIT           int64      `json:"nit"`
	RazonSocial   string     `json:"razonSocial"`
	Representante string     `json:"representante"`
	Direccion     *string    `json:"direccion,omitempty"`
	Telefono      *int       `json:"telefono,omitzero"`
	Email         *string    `json:"email,omitempty"`
	Celular       *int       `json:"celular,omitzero"`
	Estado        string     `json:"estado"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}
type ProveedorDetail struct {
	Proveedor
}

type ProveedorRequest struct {
	NIT           int64   `json:"nit"`
	RazonSocial   string  `json:"razonSocial"`
	Representante string  `json:"representante"`
	Direccion     *string `json:"direccion,omitempty"`
	Telefono      *int    `json:"telefono,omitzero"`
	Email         *string `json:"email,omitempty"`
	Celular       *int    `json:"celular,omitzero"`
}

type ProveedorSimple struct {
	Id          int    `json:"id"`
	NIT         int64  `json:"nit"`
	RazonSocial string `json:"razonSocial"`
}
type ProveedorInfo struct {
	Id            int        `json:"id"`
	NIT           int64      `json:"nit"`
	RazonSocial   string     `json:"razonSocial"`
	Representante string     `json:"representante"`
	Direccion     *string    `json:"direccion,omitempty"`
	Estado        string     `json:"estado"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}
