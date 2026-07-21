// Package sqlserver implementa el adaptador de persistencia (driven
// adapter) para SQL Server 2022 usando database/sql y el driver oficial
// microsoft/go-mssqldb.
package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

// Config agrupa los parámetros de conexión y de pool.
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewConnection abre el pool de conexiones hacia SQL Server. database/sql
// ya administra pooling internamente; aquí solo se ajustan sus límites para
// un entorno de alta concurrencia.
func NewConnection(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("opening sqlserver connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging sqlserver: %w", err)
	}

	return db, nil
}
