package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, prosseguindo com variáveis de ambiente do sistema.")
	}
}
