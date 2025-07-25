package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type ClienteInfo struct {
	Id          uint        `json:"id"`
	NitCi       *uint       `json:"nitCi"`
	Complemento pgtype.Text `json:"complemento"`
	Tipo        string      `json:"tipo"`
	Estado      string      `json:"estado"`
	RazonSocial string      `json:"razonSocial"`
}

type ClienteDetail struct {
	Id          uint       `json:"id"`
	NitCi       *uint      `json:"nitCi"`
	Complemento *string    `json:"complemento"`
	Tipo        string     `json:"tipo"`
	RazonSocial string     `json:"razonSocial"`
	Email       string     `json:"email"`
	Telefono    *uint      `json:"telefono"`
	Estado      string     `json:"estado"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type ClienteRequest struct {
	NitCi       *uint   `json:"nitCi"`
	Complemento *string `json:"complemento"`
	Tipo        string  `json:"tipo"`
	RazonSocial string  `json:"razonSocial"`
	Email       string  `json:"email"`
	Telefono    *uint   `json:"telefono"`
}

type ClienteSimple struct {
	Id          uint    `json:"id"`
	NitCi       *uint   `json:"nitCi"`
	Complemento *string `json:"complemento"`
	RazonSocial string  `json:"razonSocial"`
}

type ClienteId struct {
	Id int `json:"id"`
}
