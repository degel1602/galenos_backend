package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("sqlserver", os.Getenv("SQLSERVER_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
SELECT 
     i.name AS IndexName,
     c.name AS ColumnName
FROM 
     sys.indexes i
INNER JOIN 
     sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
INNER JOIN 
     sys.columns c ON ic.object_id = c.object_id AND c.column_id = ic.column_id
WHERE 
     i.object_id = OBJECT_ID('pacientes')
ORDER BY 
     i.name, ic.index_column_id;
	`

	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("=== Índices en la tabla 'pacientes' ===")
	for rows.Next() {
		var indexName sql.NullString
		var columnName string
		if err := rows.Scan(&indexName, &columnName); err == nil {
			idx := "HEAP/NO_NAME"
			if indexName.Valid {
				idx = indexName.String
			}
			fmt.Printf("Index: %-30s | Column: %s\n", idx, columnName)
		}
	}
}
