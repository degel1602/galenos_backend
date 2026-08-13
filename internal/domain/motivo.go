package domain

// MotivoAtencion representa un motivo de atención para una evolución
type MotivoAtencion struct {
	IdMotivo      int    `json:"idMotivo"`
	IdRegAtencion int    `json:"idRegAtencion"`
	Motivo        string `json:"motivo"`
	FechaRegistro string `json:"fechaRegistro"`
}
