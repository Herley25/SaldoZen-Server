package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	conn := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Banco conectado com sucesso")
}
