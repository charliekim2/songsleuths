package db

import (
	"database/sql"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	url := os.Getenv("DB_URL")
	token := os.Getenv("DB_AUTH")

	authUrl := url + "?authToken=" + token

	sqliteDB, err := sql.Open("libsql", authUrl)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqliteDB,
	}), &gorm.Config{})

	return db, nil
}
