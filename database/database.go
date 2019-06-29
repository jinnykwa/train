package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func GetDBConn() (*sql.DB, error) {
	//url := os.Getenv("DATABASE_URL")
	//fmt.Println("url", url)
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
