package httpadapter

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/galenos-pro/appointments-api/docs"
)

// NewRouter arma el árbol de rutas de la REST API, incluyendo CORS para que
// el frontend Angular pueda consumirla desde otro origen.
func NewRouter(
	appointmentHandler *AppointmentHandler,
	patientHandler *PatientHandler,
	catalogHandler *CatalogHandler,
	reniecHandler *ReniecHandler,
	sisHandler *SisHandler,
	triageHandler *TriageHandler,
	allowedOrigins []string,
) *gin.Engine {
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
			pacientes.GET("/buscar", patientHandler.Search)
			pacientes.GET("/por-documento", patientHandler.GetByDocumentAndType)
			pacientes.GET("/:idOrDoc", patientHandler.Get)
			pacientes.PUT("/:id", patientHandler.Update)
		}

		etnias := v1.Group("/etnias")
		{
			etnias.GET("", catalogHandler.ListEtnias)
		}

		idiomas := v1.Group("/idiomas")
		{
			idiomas.GET("", catalogHandler.ListIdiomas)
		}

		v1.GET("/tipos-sexo", catalogHandler.ListTipoSexos)
		v1.GET("/estados-civil", catalogHandler.ListEstadosCivil)
		v1.GET("/grados-instruccion", catalogHandler.ListGradosInstruccion)
		v1.GET("/ocupaciones", catalogHandler.ListOcupaciones)
		v1.GET("/tipos-documentos", catalogHandler.ListTiposDocumentos)

		v1.GET("/departamentos", catalogHandler.ListDepartamentos)
		v1.GET("/provincias/:id", catalogHandler.ListProvincias)
		v1.GET("/distritos/:id", catalogHandler.ListDistritos)
		v1.GET("/centros-poblados/:id", catalogHandler.ListCentrosPoblados)
		v1.GET("/paises", catalogHandler.ListPaises)
		v1.GET("/estados-llego-paciente", catalogHandler.ListEstadosLlegoPaciente)
		v1.GET("/fuentes-financiamiento", catalogHandler.ListFuentesFinanciamiento)
		v1.GET("/servicios/:idTipoServicio", catalogHandler.ListServicios)
		v1.GET("/datos-institucion", catalogHandler.GetDatosInstitucion)
		v1.GET("/especialidades", catalogHandler.ListEspecialidades)

		reniec := v1.Group("/reniec")
		{
			reniec.GET("/:nrodoc", reniecHandler.Consultar)
		}

		sis := v1.Group("/sis")
		{
			sis.GET("/afiliado/:nrodoc", sisHandler.ConsultarAfiliado)
			sis.POST("/filiaciones", sisHandler.GestionarAfiliacion)
			sis.POST("/fua", sisHandler.ForzarGuardadoFua)
			sis.POST("/fua/agregar", sisHandler.AgregarFua)
			sis.GET("/fua/imprimir", sisHandler.GetFuaImprimir)
			sis.GET("/diagnosticos", sisHandler.ListDiagnosticos)
			sis.GET("/medicamentos", sisHandler.ListMedicamentos)
			sis.GET("/procedimientos", sisHandler.ListProcedimientos)
			sis.GET("/consumo", sisHandler.ListConsumo)
		}

		triaje := v1.Group("/triaje")
		{
			triaje.GET("", triageHandler.List)
			triaje.GET("/pendientes-admision", triageHandler.ListPendingAdmission)
			triaje.GET("/reporte", triageHandler.GetReporte)
			triaje.GET("/ficha-admision", triageHandler.GetFichaAdmision)
			triaje.POST("", triageHandler.Create)
			triaje.POST("/admision", triageHandler.CreateAdmission)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
