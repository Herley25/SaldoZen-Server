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

// GetMonthlySummary retorna o resumo mensal de despesas e receitas de um usuário
//
// @Summary Resumo mensal
// @Tags Summary
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param month query string true "Mês (1-12)"
// @Param year query string true "Ano (YYYY)"
// @Success 200 {object} models.Summary
// @Failure 400,401,500 {string} string
// @Router /summary/{userId} [get]
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

	// Consulta para obter o resumo mensal
	query := `
		SELECT
			COALESCE(SUM(valor),0) AS total_despesas,
			COALESCE(SUM(valor),0) FILTER (WHERE paga=true)                          AS total_pagas,
			COALESCE(SUM(valor),0) FILTER (WHERE paga=false AND vencimento>NOW())    AS pendentes,
			COALESCE(SUM(valor),0) FILTER (WHERE paga=false AND vencimento<NOW())    AS total_vencidas,
			(SELECT COALESCE(SUM(valor),0)
			 FROM incomes
			 WHERE user_id=$1 AND data_recebimento >= $2 AND data_recebimento < $3) AS receitas
		FROM expenses
		WHERE user_id=$1 AND vencimento >= $2 AND vencimento < $3;
	`

	err = db.DB.QueryRow(query, userID, startDate, endDate).Scan(
		&summary.TotalPagas,
		&summary.Pendentes,
		&summary.TotalVencidas,
		&summary.TotalDespesas,
		&summary.Receitas,
	)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Erro ao buscar resumo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	summary.Saldo = summary.Receitas - summary.TotalDespesas // Calcular o saldo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)
}
