package domain

type Persona struct {
	Id              int32   `json:"id"`
	Ci              int32   `json:"ci"`
	Complemento     *string `json:"complemento"`
	Nombres         string  `json:"nombres"`
	ApellidoPaterno string  `json:"apellidoPaterno"`
	ApellidoMaterno string  `json:"apellidoMaterno"`
	Genero          string  `json:"genero"`
}

type PersonaRequest struct {
	Ci              int32   `json:"ci"`
	Complemento     *string `json:"complemento"`
	Nombres         string  `json:"nombres"`
	ApellidoPaterno string  `json:"apellidoPaterno"`
	ApellidoMaterno string  `json:"apellidoMaterno"`
	Genero          string  `json:"genero"`
}
