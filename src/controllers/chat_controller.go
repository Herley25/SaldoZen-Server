package controllers

import (
	"encoding/json"
	"finance/src/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CategoryChart struct {
	Categoria string  `json:"categoria"`
	Total     float64 `json:"total"`
}

type StatusChart struct {
	Status string  `json:"status"`
	Total  float64 `json:"total"`
}

type MonthChart struct {
	Mes   string  `json:"mes"`
	Total float64 `json:"total"`
}

// GetExpensesByCategory retorna soma das despesas agrupadas por categoria
//
// @Summary Despesas por categoria
// @Tags Charts
// @Param userId path string true "ID do usuário"
// @Param month query string true "Mês (1-12)"
// @Param year query string true "Ano (YYYY)"
// @Success 200 {array} CategoryChart
// @Failure 400,401,500 {string} string
// @Router /charts/expenses-by-category/{userId} [get]
func GetExpensesByCategory(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["userId"]

	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, err1 := strconv.Atoi(monthStr)
	year, err2 := strconv.Atoi(yearStr)
	if err1 != nil || err2 != nil || month < 1 || month > 12 {
		http.Error(w, "Parâmetros de mês/ano inválidos", http.StatusBadRequest)
		return
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC) // Define o início do mês
	end := start.AddDate(0, 1, 0)                                        // Adiciona um mês para definir o final do mês

	rows, err := db.DB.Query(`
		SELECT categoria, COALESCE(SUM(valor), 0) AS total
		FROM expenses
		WHERE user_id = $1 AND vencimento >= $2 AND vencimento < $3
		GROUP BY categoria
		ORDER BY total DESC
	`, userID, start, end)
	if err != nil {
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Fecha as linhas após o uso

	var result []CategoryChart
	for rows.Next() {
		var row CategoryChart
		if err := rows.Scan(&row.Categoria, &row.Total); err != nil {
			http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
			return
		}
		result = append(result, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetExpensesByStatus retorna soma das despesas agrupadas por status no mês
//
// @Summary Despesas por status
// @Tags Charts
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param month query string true "Mês (1-12)"
// @Param year query string true "Ano (YYYY)"
// @Success 200 {array} StatusChart
// @Failure 400,401,500 {string} string
// @Router /charts/expenses-by-status/{userId} [get]
func GetExpensesByStatus(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["userId"]
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if month < 1 || month > 12 || year == 0 {
		http.Error(w, "Parâmetros de mês/ano inválidos", http.StatusBadRequest)
		return
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC) // Define o início do mês
	end := start.AddDate(0, 1, 0)                                        // Adiciona um mês para definir o final do mês

	query := `
		SELECT
			CASE
				WHEN paga = true THEN 'Paga'
				WHEN paga = false AND vencimento < CURRENT_DATE THEN 'Vencida'
				ELSE 'A Vencer'
			END AS status,
			COALESCE(SUM(valor), 0) AS total
		FROM expenses
		WHERE user_id = $1 AND vencimento >= $2 AND vencimento < $3
		GROUP BY status
		ORDER BY status
	`
	rows, err := db.DB.Query(query, uid, start, end)
	if err != nil {
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Fecha as linhas após o uso

	var result []StatusChart
	for rows.Next() {
		var s StatusChart
		if err := rows.Scan(&s.Status, &s.Total); err != nil {
			http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
			return
		}
		result = append(result, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetMonthlySummaryChart retorna totais mensais de despesas por mês
//
// @Summary Resumo mensal do ano
// @Tags Charts
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param year query string true "Ano (YYYY)"
// @Success 200 {array} MonthChart
// @Failure 400,401,500 {string} string
// @Router /charts/monthly-summary/{userId} [get]
func GetMonthlySummaryChart(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["userId"]
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		http.Error(w, "Ano obrigatório", http.StatusBadRequest)
		return
	}

	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC) // Define o início do ano
	end := start.AddDate(1, 0, 0)                        // Adiciona um ano para definir o final do ano

	query := `
		SELECT EXTRACT(MONTH FROM vencimento)::INT As mes, COALESCE(SUM(valor), 0) AS total
		FROM expenses
		WHERE user_id = $1 AND vencimento >= $2 AND vencimento < $3
		GROUP BY mes
		ORDER BY mes
	`
	rows, err := db.DB.Query(query, uid, start, end)
	if err != nil {
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Fecha as linhas após o uso

	var result []MonthChart
	for rows.Next() {
		var m MonthChart
		if err := rows.Scan(&m.Mes, &m.Total); err != nil {
			http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
			return
		}
		result = append(result, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
