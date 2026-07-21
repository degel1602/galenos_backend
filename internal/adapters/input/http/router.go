package httpadapter

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter arma el árbol de rutas de la REST API, incluyendo CORS para que
// el frontend Angular pueda consumirla desde otro origen.
func NewRouter(appointmentHandler *AppointmentHandler, patientHandler *PatientHandler, allowedOrigins []string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/api/v1")
	{
		appointments := v1.Group("/appointments")
		{
			appointments.POST("", appointmentHandler.Create)
			appointments.GET("/:id", appointmentHandler.GetByID)
		}

		pacientes := v1.Group("/pacientes")
		{
			pacientes.GET("", patientHandler.List)
			pacientes.GET("/:numDocumento", patientHandler.GetByDocumentNumber)
		}
	}

	return router
}
