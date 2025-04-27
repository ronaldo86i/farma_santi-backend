package domain

import "time"

type Proveedor struct {
	Id            int        `json:"id"`
	NIT           int        `json:"nit"`
	Nombre        string     `json:"nombre"`
	Representante string     `json:"representante"`
	Direccion     string     `json:"direccion"`
	Telefono      *int       `json:"telefono"`
	Email         *string    `json:"email,"`
	Celular       *int       `json:"celular"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}
type ProveedorDetail struct {
	Proveedor
}

type ProveedorRequest struct {
	NIT           int     `json:"nit"`
	Nombre        string  `json:"nombre"`
	Representante string  `json:"representante"`
	Direccion     string  `json:"direccion"`
	Telefono      *int    `json:"telefono"`
	Email         *string `json:"email"`
	Celular       *int    `json:"celular"`
}

type ProveedorInfo struct {
	Id            int        `json:"id"`
	NIT           int        `json:"nit"`
	Nombre        string     `json:"nombre"`
	Representante string     `json:"representante"`
	Direccion     string     `json:"direccion"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}
