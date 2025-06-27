package domain

import "time"

type ClienteInfo struct {
	Id          uint   `json:"id"`
	NitCi       *uint  `json:"nitCi"`
	Complemento string `json:"complemento"`
	Tipo        string `json:"tipo"`
	Estado      string `json:"estado"`
	RazonSocial string `json:"razonSocial"`
}

type ClienteDetail struct {
	Id          uint       `json:"id"`
	NitCi       *uint      `json:"nitCi"`
	Complemento string     `json:"complemento"`
	Tipo        string     `json:"tipo"`
	RazonSocial string     `json:"razonSocial"`
	Email       string     `json:"email"`
	Telefono    uint       `json:"telefono"`
	Estado      string     `json:"estado"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type ClienteRequest struct {
	NitCi       *uint  `json:"nitCi"`
	Complemento string `json:"complemento"`
	Tipo        string `json:"tipo"`
	RazonSocial string `json:"razonSocial"`
	Email       string `json:"email"`
	Telefono    uint   `json:"telefono"`
}
