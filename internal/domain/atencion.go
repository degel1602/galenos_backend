package domain

// AdmisionDesdeTriaje agrupa los datos necesarios para crear la atención
// (admisión) de un paciente a partir de su triaje, invocando el SP
// WebCrearAtencionDesdeTriaje. Los campos opcionales son punteros para
// distinguir "no enviado" (NULL) de un valor cero.
type AdmisionDesdeTriaje struct {
	IDTriaje            *int64  `json:"idTriaje"`
	IDPacienteTriaje    *int64  `json:"idPacienteTriaje"`
	IDEmpleado          *int64  `json:"idEmpleado"`
	NombreAcompanante   *string `json:"nombreAcompanante"`
	TelefonoAcompanante *string `json:"telefonoAcompanante"`
	DireccionPaciente   *string `json:"direccionPaciente"`
	Observacion         *string `json:"observacion"`
}
