package domain

// ReniecResult es la respuesta normalizada de una consulta al servicio web
// RENIEC. Resultado contiene los valores crudos que devuelve el SOAP (según la
// operación: basico o completo) y Datos agrupa los campos ya interpretados con
// un layout estable para que el consumidor no dependa de índices del arreglo.
type ReniecResult struct {
	DNI       string
	Operacion string
	Resultado []string
	Datos     ReniecDatos
}

// ReniecDatos agrupa los datos de la persona extraídos de la respuesta de
// RENIEC con un layout estable (independiente de la operación SOAP). Sexo,
// estado civil y domicilio quedan vacíos si el servicio no los devuelve o no
// pueden inferirse con certeza.
type ReniecDatos struct {
	ApellidoPaterno        string
	ApellidoMaterno        string
	Nombres                string
	PrimerNombre           string
	SegundoNombre          string
	TercerNombre           string
	FechaNacimiento        string
	Sexo                   string
	EstadoCivil            string
	Departamento           string
	Provincia              string
	Distrito               string
	Direccion              string
	Ubigeo                 string
	NombrePadre            string
	NombreMadre            string
	DepartamentoNacimiento string
	ProvinciaNacimiento    string
	DistritoNacimiento     string
}
