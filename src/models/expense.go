package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Descricao     string     `json:"descricao"`
	Valor         float64    `json:"valor"`
	Vencimento    time.Time  `json:"vencimento"`
	Paga          bool       `json:"paga"`
	DataPagamento *time.Time `json:"data_pagamento,omitempty"`
	Categoria     string     `json:"categoria"`
	Observacoes   *string    `json:"observacoes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	Status        string     `json:"status"`
}

func (e *Expense) StatusHoje() string {
	hoje := time.Now()
	if e.Paga {
		return "Paga"
	}
	if hoje.After(e.Vencimento) {
		return "Vencida"
	}
	return "A Vencer"
}
