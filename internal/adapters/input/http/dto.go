package httpadapter

import "time"

// createAppointmentRequest es el contrato de entrada esperado desde Angular.
type createAppointmentRequest struct {
	PatientID string    `json:"patientId" binding:"required"`
	DoctorID  string    `json:"doctorId" binding:"required"`
	StartsAt  time.Time `json:"startsAt" binding:"required"`
	EndsAt    time.Time `json:"endsAt" binding:"required"`
	Reason    string    `json:"reason" binding:"required,max=500"`
}
