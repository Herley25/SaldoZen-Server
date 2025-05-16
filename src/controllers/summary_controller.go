package controllers

import (
	"database/sql"
	"encoding/json"
	"finance/src/db"
	"finance/src/models"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}

	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		http.Error(w, "Mês inválido", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1 {
		http.Error(w, "Ano inválido", http.StatusBadRequest)
		return
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	var summary models.Summary
	summary.Mes = month
	summary.Ano = year

	query := `
		SELECT
			COALESCE(SUM(valor), 0) FILTER (WHERE paga = true) as total_pagas,
			COALESCE(SUM(valor), 0) FILTER (WHERE paga = false AND vencimento > CURRENT_DATE) as pendentes,
			COALESCE(SUM(valor), 0) FILTER (WHERE paga = false AND vencimento < CURRENT_DATE) as total_vencidas,
			COALESCE(SUM(valor), 0) as total_despesas
		FROM expenses
		WHERE user_id = $1 AND vencimento >= $2 AND vencimento < $3
	`

	err = db.DB.QueryRow(query, userID, startDate, endDate).Scan(
		&summary.TotalPagas,
		&summary.Pendentes,
		&summary.TotalVencidas,
		&summary.TotalDespesas,
	)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Erro ao buscar resumo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// saldo pode ser alterado futuramente com base em receitas, mas aqui será apenas o total de despesas negativas
	summary.Saldo = -summary.TotalDespesas

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}
