package main

import (
	"log"
	"net/http"

	"finance/src/config"
	"finance/src/db"
	"finance/src/routes"
)

func main() {
	config.LoadEnv()
	db.Init()
	router := routes.SetupRoutes()
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
