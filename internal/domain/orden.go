package domain

type OrdenMedica struct {
	IdOrden       int            `json:"idOrden"`
	IdRegAtencion int            `json:"idRegAtencion"`
	IdMedico      int            `json:"idMedico"`
	Medico        string         `json:"medico"`
	FechaOrden    string         `json:"fechaOrden"`
	Estado        string         `json:"estado"`
	Observacion   string         `json:"observacion"`
	Detalles      []DetalleOrden `json:"detalles"`
}

type DetalleOrden struct {
	IdProducto     int     `json:"idProducto"`
	NombreProducto string  `json:"nombreProducto"`
	Codigo         string  `json:"codigo"`
	Cantidad       int     `json:"cantidad"`
	Precio         float64 `json:"precio"`
	Total          float64 `json:"total"`
	Indicaciones   string  `json:"indicaciones"`
}

type ProductoCatalogo struct {
	IdProducto        int     `json:"idProducto"`
	Codigo            string  `json:"codigo"`
	Nombre            string  `json:"nombre"`
	Concentracion     string  `json:"concentracion"`
	Presentacion      string  `json:"presentacion"`
	FormaFarmaceutica string  `json:"formaFarmaceutica"`
	PrecioVenta       float64 `json:"precioVenta"`
}
