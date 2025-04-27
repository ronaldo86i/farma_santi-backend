package domain

type Persona struct {
	Id              uint    `json:"id"`
	Ci              int     `json:"ci"`
	Complemento     *string `json:"complemento"`
	Nombres         string  `json:"nombres"`
	ApellidoPaterno string  `json:"apellidoPaterno"`
	ApellidoMaterno string  `json:"apellidoMaterno"`
	Genero          string  `json:"genero"`
}

type PersonaRequest struct {
	Ci              int     `json:"ci"`
	Complemento     *string `json:"complemento"`
	Nombres         string  `json:"nombres"`
	ApellidoPaterno string  `json:"apellidoPaterno"`
	ApellidoMaterno string  `json:"apellidoMaterno"`
	Genero          string  `json:"genero"`
}
