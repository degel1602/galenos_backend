package domain

// ReniecResult es la respuesta normalizada de una consulta al servicio web
// RENIEC. Resultado contiene los valores crudos que devuelve el SOAP (según la
// operación: basico o completo) y Datos agrupa los campos ya interpretados con
// un layout estable para que el consumidor no dependa de índices del arreglo.
type ReniecResult struct {
	DNI       string      `json:"dni"`
	Operacion string      `json:"operacion"`
	Resultado []string    `json:"resultado"`
	Datos     ReniecDatos `json:"datos"`
}

// ReniecDatos agrupa los datos de la persona extraídos de la respuesta de
// RENIEC con un layout estable (independiente de la operación SOAP). Sexo
// queda vacío si el servicio no lo devuelve o no puede inferirse.
type ReniecDatos struct {
	ApellidoPaterno string `json:"apellidoPaterno"`
	ApellidoMaterno string `json:"apellidoMaterno"`
	Nombres         string `json:"nombres"`
	PrimerNombre    string `json:"primerNombre"`
	SegundoNombre   string `json:"segundoNombre"`
	TercerNombre    string `json:"tercerNombre"`
	FechaNacimiento string `json:"fechaNacimiento"`
	Sexo            string `json:"sexo"`
}
