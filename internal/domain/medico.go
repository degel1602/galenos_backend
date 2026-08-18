package domain

// MedicoFila representa un registro retornado por el SP
// usp_go_MedicosFiltrarPorIdEspecialidad. El SP devuelve las columnas
// IdMedico, ApellidoPaterno, ApellidoMaterno y Nombre; el NombreCompleto
// se arma concatenando las partes.
type MedicoFila struct {
	IdMedico        int     `json:"idMedico"`
	CodigoPlanilla  *string `json:"codigoPlanilla"`
	ApellidoPaterno *string `json:"apellidoPaterno"`
	ApellidoMaterno *string `json:"apellidoMaterno"`
	Nombre          *string `json:"nombre"`
	Especialidad    *string `json:"especialidad"`
	Colegiatura     *string `json:"colegiatura"`
	RNE             *string `json:"rne"`
	NombreCompleto  string  `json:"nombreCompleto"`
}
