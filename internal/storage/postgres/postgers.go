package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(url string) (*sql.DB, error) {
	return sql.Open("pgx", url)
}
