package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/galenos-pro/appointments-api/internal/config"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlserver", cfg.SQLServerDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var definition string
	err = db.QueryRowContext(context.Background(), "SELECT OBJECT_DEFINITION(OBJECT_ID('usp_go_Login'))").Scan(&definition)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SP Definition:")
	fmt.Println(definition)
}
