//	@title			SaldoZen API
//	@version		1.0
//	@description	API de controle financeiro pessoal.

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

package main

import (
	"log"
	"net/http"

	"finance/src/config"
	"finance/src/db"
	"finance/src/routes"

	_ "finance/src/controllers"
	_ "finance/src/docs" // Importando os docs gerados pelo Swag
	_ "finance/src/models"

	_ "github.com/swaggo/http-swagger" // Importando o Swagger para documentação
)

func main() {
	config.LoadEnv()
	db.Init()
	router := routes.SetupRoutes()
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
