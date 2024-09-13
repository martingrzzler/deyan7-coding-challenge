package persist

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=127.0.0.1 port=5432 dbname=postgres user=postgres password=password sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return db, nil
}
