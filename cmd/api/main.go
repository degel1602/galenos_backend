// Command api arranca la REST API del sistema hospitalario. Es el único
// punto donde se conectan las implementaciones concretas de los adaptadores
// con los puertos que el dominio define.
package main

// @title Galenos Pro Appointments API
// @version 1.0
// @description REST API del sistema hospitalario Galenos Pro. Arquitectura hexagonal, persiste en SQL Server 2022.
// @host localhost:8080
// @BasePath /api/v1

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/galenos-pro/appointments-api/docs"
	httpadapter "github.com/galenos-pro/appointments-api/internal/adapters/input/http"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/persistence/sqlserver"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/reniec"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/sis"
	"github.com/galenos-pro/appointments-api/internal/config"
	"github.com/galenos-pro/appointments-api/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}

func run() error {
	// Carga .env solo si existe (desarrollo local); en producción las
	// variables de entorno se inyectan directamente en el proceso.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// --- Adaptador de salida: persistencia SQL Server ---
	db, err := sqlserver.NewConnection(sqlserver.Config{
		DSN:             cfg.SQLServerDSN,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	appointmentRepo := sqlserver.NewAppointmentRepository(db)
	patientRepo := sqlserver.NewPatientRepository(db)
	catalogRepo := sqlserver.NewCatalogRepository(db)
	triageRepo := sqlserver.NewTriageRepository(db)
	sisRepo := sqlserver.NewSisRepository(db)

	// --- Adaptador de salida: servicio externo RENIEC ---
	reniecClient := reniec.New(reniec.Config{
		App:     cfg.ReniecApp,
		Usuario: cfg.ReniecUsuario,
		Clave:   cfg.ReniecClave,
		URL:     cfg.ReniecURL,
		Timeout: cfg.ReniecTimeout,
	})

	// --- Adaptador de salida: servicio externo SIS ---
	sisClient := sis.New(sis.Config{
		Usuario:       cfg.SISUsuario,
		Clave:         cfg.SISClave,
		URL:           cfg.SISURL,
		DNIAutorizado: cfg.SISDNIAutorizado,
		Timeout:       cfg.SISTimeout,
	})

	// --- Núcleo de dominio: casos de uso implementando los puertos de entrada ---
	appointmentService := usecase.NewAppointmentUseCase(appointmentRepo)
	patientService := usecase.NewPatientUseCase(patientRepo)
	catalogService := usecase.NewCatalogUseCase(catalogRepo)
	reniecService := usecase.NewReniecUseCase(reniecClient)
	sisService := usecase.NewSisUseCase(sisClient, sisRepo)
	triageService := usecase.NewTriageUseCase(triageRepo)

	// El host de Swagger se ajusta en runtime al IP detectado de la
	// máquina (o al configurado en SERVER_HOST) para que otros equipos de
	// la red consuman la API por IP en lugar de localhost.
	docs.SwaggerInfo.Host = cfg.ServerHost + ":" + cfg.ServerPort

	// --- Adaptador de entrada: HTTP/REST para Angular ---
	appointmentHandler := httpadapter.NewAppointmentHandler(appointmentService)
	patientHandler := httpadapter.NewPatientHandler(patientService)
	catalogHandler := httpadapter.NewCatalogHandler(catalogService)
	reniecHandler := httpadapter.NewReniecHandler(reniecService)
	sisHandler := httpadapter.NewSisHandler(sisService)
	triageHandler := httpadapter.NewTriageHandler(triageService)
	router := httpadapter.NewRouter(
		appointmentHandler,
		patientHandler,
		catalogHandler,
		reniecHandler,
		sisHandler,
		triageHandler,
		cfg.AllowedOrigins,
	)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on http://%s:%s (swagger: http://%s:%s/swagger/index.html)", cfg.ServerHost, cfg.ServerPort, cfg.ServerHost, cfg.ServerPort)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on :%s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
