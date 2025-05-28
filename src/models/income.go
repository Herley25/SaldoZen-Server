package models

import (
	"time"

	"github.com/google/uuid"
)

type Income struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Descricao       string    `json:"descricao"`
	Valor           float64   `json:"valor"`
	DataRecebimento time.Time `json:"data_recebimento"`
	Categoria       string    `json:"categoria"`
	Observacoes     *string   `json:"observacoes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
