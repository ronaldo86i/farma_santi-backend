package domain

import "time"

type Proveedor struct {
	Id            int32      `json:"id"`
	NIT           int32      `json:"nit"`
	RazonSocial   string     `json:"razonSocial"`
	Representante string     `json:"representante"`
	Direccion     *string    `json:"direccion,omitempty"`
	Telefono      *int32     `json:"telefono,omitzero"`
	Email         *string    `json:"email,omitempty"`
	Celular       *int32     `json:"celular,omitzero"`
	Estado        string     `json:"estado"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}
type ProveedorDetail struct {
	Proveedor
}

type ProveedorRequest struct {
	NIT           int32   `json:"nit"`
	RazonSocial   string  `json:"razonSocial"`
	Representante string  `json:"representante"`
	Direccion     *string `json:"direccion,omitempty"`
	Telefono      *int32  `json:"telefono,omitzero"`
	Email         *string `json:"email,omitempty"`
	Celular       *int32  `json:"celular,omitzero"`
}

type ProveedorInfo struct {
	Id            int32      `json:"id"`
	NIT           int32      `json:"nit"`
	RazonSocial   string     `json:"razonSocial"`
	Representante string     `json:"representante"`
	Direccion     *string    `json:"direccion,omitempty"`
	Estado        string     `json:"estado"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}
