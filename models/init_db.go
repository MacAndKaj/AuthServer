package models

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type UsersDatabase struct {
	dbHandle *sql.DB
	logger   *log.Logger
}

func InitUsersDatabase(l *log.Logger) *UsersDatabase {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "root"
	cfg.DBName = "users_db"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	l.Println("Database configured successfully")

	return &UsersDatabase{
		dbHandle: db,
		logger:   l,
	}
}

func (db *UsersDatabase) Shutdown() error {
	db.logger.Println("Shutdown database")

	return db.dbHandle.Close()
}
