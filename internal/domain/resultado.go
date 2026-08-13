package domain

type Resultado struct {
	IdResultado   int    `json:"idResultado"`
	IdPaciente    int    `json:"idPaciente"`
	TipoResultado string `json:"tipoResultado"` // "Laboratorio" or "Imagen"
	NombreExamen  string `json:"nombreExamen"`
	FechaExamen   string `json:"fechaExamen"`
	Detalle       string `json:"detalle"`
	Estado        string `json:"estado"`
}
