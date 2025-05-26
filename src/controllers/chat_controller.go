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
