package domain

type OrdenMedica struct {
	IdOrden       int    `json:"idOrden"`
	IdRegAtencion int    `json:"idRegAtencion"`
	IdMedico      int    `json:"idMedico"`
	FechaOrden    string `json:"fechaOrden"`
	Estado        string `json:"estado"`
	Observacion   string `json:"observacion"`
}

type DetalleOrden struct {
	IdDetalleOrden int    `json:"idDetalleOrden"`
	IdOrden        int    `json:"idOrden"`
	IdServicio     int    `json:"idServicio"`
	NombreServicio string `json:"nombreServicio"`
	Cantidad       int    `json:"cantidad"`
}
