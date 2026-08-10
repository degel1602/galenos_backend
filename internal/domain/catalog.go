package domain

// Etnia es un registro del catálogo de etnias (SP ups_go_ListarEtnias).
type Etnia struct {
	Codigo      string  `json:"codigo"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// Idioma es un registro del catálogo de lenguas (SP ups_go_ListarIdiomas).
type Idioma struct {
	ID     int64   `json:"id"`
	Lengua *string `json:"lengua,omitempty"`
}

// TipoSexo es un registro del catálogo de sexos (SP usp_go_ListarTiposSexos).
type TipoSexo struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// TipoEstadoCivil es un registro del catálogo de estados civiles
// (SP usp_go_ListarEstadosCivil).
type TipoEstadoCivil struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// TipoGradoInstruccion es un registro del catálogo de grados de instrucción
// (SP usp_go_ListarGradoInstruccion).
type TipoGradoInstruccion struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// TipoOcupacion es un registro del catálogo de ocupaciones
// (SP usp_go_ListarOcupaciones).
type TipoOcupacion struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// TipoDocumento es un registro del catálogo de tipos de documento de
// identidad (SP usp_go_ListarTiposDocumentos).
type TipoDocumento struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// Departamento es un registro de la tabla Departamentos
// (SP usp_go_ListarDepartamentos).
type Departamento struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// Provincia es un registro de la tabla Provincias
// (SP usp_go_ListarProvincias @IdDepartamento).
type Provincia struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// Distrito es un registro de la tabla Distritos
// (SP usp_go_ListarDistritos @IdProvincia).
type Distrito struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// CentroPoblado es un registro del catálogo de centros poblados
// (SP usp_go_ListarCentrosPoblados @IdDistrito).
type CentroPoblado struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// Pais es un registro del catálogo de países (SP usp_go_ListarPaises).
type Pais struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// EstadoLlegoPaciente es un registro del catálogo de estados de llegada del
// paciente (SP usp_go_listarEstadosLlegoPaciente).
type EstadoLlegoPaciente struct {
	ID          int64   `json:"id"`
	Descripcion *string `json:"descripcion,omitempty"`
}

// FuenteFinanciamiento es un registro del catálogo de fuentes de
// financiamiento (SP usp_go_ListarFuentesFinanciamiento).
type FuenteFinanciamiento struct {
	ID                   int64   `json:"idFuenteFinanciamiento"`
	Descripcion          *string `json:"descripcion,omitempty"`
	IdTipoFinanciamiento int64   `json:"idTipoFinanciamiento"`
}

// Servicio es un registro del catálogo de servicios por tipo
// (SP usp_go_ListarServicios @IdTipoServicio).
type Servicio struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}

// DatosInstitucion contiene los datos de la institución (EE.SS.) que el
// SP webParametrosDatosInstitucion devuelve en una única fila.
type DatosInstitucion struct {
	RucEESS    *string `json:"rucEess,omitempty"`
	Nombre     *string `json:"nombre,omitempty"`
	Direccion  *string `json:"direccion,omitempty"`
	Telefono   *string `json:"telefono,omitempty"`
	Codigo     *string `json:"codigo,omitempty"`
	CodRenaes  *string `json:"codRenaes,omitempty"`
	LogoHospi  *string `json:"logoHospi,omitempty"`
	LogoMinsa  *string `json:"logoMinsa,omitempty"`
	UbigeoHosp *string `json:"ubigeoHosp,omitempty"`
}

// Especialidad es un registro del catálogo de especialidades
// (SP usp_go_ListarEspecialidades).
type Especialidad struct {
	ID     int64   `json:"id"`
	Nombre *string `json:"nombre,omitempty"`
}
