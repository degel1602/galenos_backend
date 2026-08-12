package httpadapter

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/galenos-pro/appointments-api/docs"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

type RouterParams struct {
	AppointmentHandler *AppointmentHandler
	PatientHandler     *PatientHandler
	CatalogHandler     *CatalogHandler
	ReniecHandler      *ReniecHandler
	SisHandler         *SisHandler
	TriageHandler      *TriageHandler
	AuthHandler        *AuthHandler
	EvolucionHandler   *EvolucionHandler
	AuthService        input.AuthService
	AllowedOrigins     []string
}

// NewRouter arma el árbol de rutas de la REST API, incluyendo CORS para que
// el frontend Angular pueda consumirla desde otro origen.
func NewRouter(p RouterParams) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     p.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", p.AuthHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(RequireBearerToken(p.AuthService))
		{
			authProtected := protected.Group("/auth")
			authProtected.GET("/menus", p.AuthHandler.GetMenus)

			pacientes := protected.Group("/pacientes")
			pacientes.GET("", p.PatientHandler.List)
			pacientes.GET("/buscar", p.PatientHandler.Search)
			pacientes.GET("/por-documento", p.PatientHandler.GetByDocumentAndType)
			pacientes.GET("/:idOrDoc", p.PatientHandler.Get)
			pacientes.PUT("/:id", p.PatientHandler.Update)
		}

		etnias := v1.Group("/etnias")
		{
			etnias.GET("", p.CatalogHandler.ListEtnias)
		}

		idiomas := v1.Group("/idiomas")
		{
			idiomas.GET("", p.CatalogHandler.ListIdiomas)
		}

		v1.GET("/tipos-sexo", p.CatalogHandler.ListTipoSexos)
		v1.GET("/estados-civil", p.CatalogHandler.ListEstadosCivil)
		v1.GET("/grados-instruccion", p.CatalogHandler.ListGradosInstruccion)
		v1.GET("/ocupaciones", p.CatalogHandler.ListOcupaciones)
		v1.GET("/tipos-documentos", p.CatalogHandler.ListTiposDocumentos)

		v1.GET("/departamentos", p.CatalogHandler.ListDepartamentos)
		v1.GET("/provincias/:id", p.CatalogHandler.ListProvincias)
		v1.GET("/distritos/:id", p.CatalogHandler.ListDistritos)
		v1.GET("/centros-poblados/:id", p.CatalogHandler.ListCentrosPoblados)
		v1.GET("/paises", p.CatalogHandler.ListPaises)
		v1.GET("/estados-llego-paciente", p.CatalogHandler.ListEstadosLlegoPaciente)
		v1.GET("/fuentes-financiamiento", p.CatalogHandler.ListFuentesFinanciamiento)
		v1.GET("/servicios/:idTipoServicio", p.CatalogHandler.ListServicios)
		v1.GET("/datos-institucion", p.CatalogHandler.GetDatosInstitucion)
		v1.GET("/especialidades", p.CatalogHandler.ListEspecialidades)

		reniec := v1.Group("/reniec")
		{
			reniec.GET("/:nrodoc", p.ReniecHandler.Consultar)
		}

		sis := v1.Group("/sis")
		{
			sis.GET("/afiliado/:nrodoc", p.SisHandler.ConsultarAfiliado)
			sis.POST("/filiaciones", p.SisHandler.GestionarAfiliacion)
			sis.POST("/fua", p.SisHandler.ForzarGuardadoFua)
			sis.POST("/fua/agregar", p.SisHandler.AgregarFua)
			sis.GET("/fua/imprimir", p.SisHandler.GetFuaImprimir)
			sis.GET("/diagnosticos", p.SisHandler.ListDiagnosticos)
			sis.GET("/medicamentos", p.SisHandler.ListMedicamentos)
			sis.GET("/procedimientos", p.SisHandler.ListProcedimientos)
			sis.GET("/consumo", p.SisHandler.ListConsumo)
		}

		triaje := v1.Group("/triaje")
		{
			triaje.GET("", p.TriageHandler.List)
			triaje.GET("/pendientes-admision", p.TriageHandler.ListPendingAdmission)
			triaje.GET("/reporte", p.TriageHandler.GetReporte)
			triaje.GET("/ficha-admision", p.TriageHandler.GetFichaAdmision)
			triaje.POST("", p.TriageHandler.Create)
			triaje.POST("/admision", p.TriageHandler.CreateAdmission)
		}

		evoluciones := v1.Group("/evoluciones")
		{
			evoluciones.POST("", p.EvolucionHandler.HandleCreateEvolucion)
			evoluciones.GET("/pacientes", p.EvolucionHandler.HandleListPatients)
			evoluciones.GET("/paciente/:pacienteId", p.EvolucionHandler.HandleListEvoluciones)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
