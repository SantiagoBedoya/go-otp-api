package mysql

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_CONN"))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
