package httpadapter

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// NewRouter arma el árbol de rutas de la REST API, incluyendo CORS para que
// el frontend Angular pueda consumirla desde otro origen. Todas las rutas
// de negocio requieren un Bearer token JWT (obtenido vía POST
// /api/v1/auth/login), excepto el propio endpoint de login.
func NewRouter(appointmentHandler *AppointmentHandler, patientHandler *PatientHandler, authHandler *AuthHandler, authService input.AuthService, allowedOrigins []string) *gin.Engine {
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
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(RequireBearerToken(authService))
		{
			appointments := protected.Group("/appointments")
			{
				appointments.POST("", appointmentHandler.Create)
				appointments.GET("/:id", appointmentHandler.GetByID)
			}

			pacientes := protected.Group("/pacientes")
			{
				pacientes.GET("", patientHandler.List)
				pacientes.GET("/:numDocumento", patientHandler.GetByDocumentNumber)
			}
		}
	}

	return router
}
