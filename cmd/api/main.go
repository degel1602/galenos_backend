// Command api arranca la REST API del sistema hospitalario. Es el único
// punto donde se conectan las implementaciones concretas de los adaptadores
// con los puertos que el dominio define.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	httpadapter "github.com/galenos-pro/appointments-api/internal/adapters/input/http"
	"github.com/galenos-pro/appointments-api/internal/adapters/output/persistence/sqlserver"
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

	// --- Núcleo de dominio: casos de uso implementando los puertos de entrada ---
	appointmentService := usecase.NewAppointmentUseCase(appointmentRepo)
	patientService := usecase.NewPatientUseCase(patientRepo)
	authService := usecase.NewAuthUseCase(cfg.APIUsername, cfg.APIPassword, cfg.JWTSecret, cfg.JWTExpiration)

	// --- Adaptador de entrada: HTTP/REST para Angular ---
	appointmentHandler := httpadapter.NewAppointmentHandler(appointmentService)
	patientHandler := httpadapter.NewPatientHandler(patientService)
	authHandler := httpadapter.NewAuthHandler(authService)
	router := httpadapter.NewRouter(appointmentHandler, patientHandler, authHandler, authService, cfg.AllowedOrigins)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

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
