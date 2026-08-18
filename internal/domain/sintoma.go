package domain

// SintomaCatalogo representa un síntoma del catálogo agrupado por sistema
type SintomaCatalogo struct {
	IdSintoma int    `json:"idSintoma"`
	Sistema   string `json:"sistema"`
	Sintoma   string `json:"sintoma"`
	Orden     int    `json:"orden"`
}

// SintomaSeleccionado representa un síntoma marcado en la evolución
type SintomaSeleccionado struct {
	IdSintoma int    `json:"idSintoma"`
	Sistema   string `json:"sistema"`
	Sintoma   string `json:"sintoma"`
}
