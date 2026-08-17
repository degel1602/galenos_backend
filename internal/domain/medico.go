package domain

// MedicoFila representa un registro retornado por el SP
// usp_go_MedicosFiltrarPorIdEspecialidad. El SP devuelve las columnas
// IdMedico, ApellidoPaterno, ApellidoMaterno y Nombre; el NombreCompleto
// se arma concatenando las partes.
type MedicoFila struct {
	IdMedico       int     `json:"idMedico"`
	ApellidoPaterno *string `json:"apellidoPaterno"`
	ApellidoMaterno *string `json:"apellidoMaterno"`
	Nombre         *string `json:"nombre"`
	NombreCompleto string  `json:"nombreCompleto"`
}
