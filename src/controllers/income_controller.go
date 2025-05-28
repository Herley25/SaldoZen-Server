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

// CreateIncome cria uma nova receita para um usuário
//
// @Summary Criar receita
// @Tags Incomes
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param income body models.Income true "Dados da receita"
// @Success 201 {object} models.Income
// @Failure 400,401,500 {string} string
// @Router /incomes/{userId} [post]
func CreateIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}

	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	income.ID = uuid.New()
	income.UserID = userID
	income.CreatedAt = time.Now()

	query := `
		INSERT INTO incomes (id, user_id, descricao, valor, data_recebimento, categoria, observacoes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = db.DB.Exec(query,
		income.ID,
		income.UserID,
		income.Descricao,
		income.Valor,
		income.DataRecebimento,
		income.Categoria,
		income.Observacoes,
		income.CreatedAt,
	)

	if err != nil {
		http.Error(w, "Erro ao criar receita: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(income)
}

// ListIncomes lista todas as receitas de um usuário em um mês e ano específicos
//
// @Summary Listar receitas
// @Tags Incomes
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param month query string true "Mês (1-12)"
// @Param year query string true "Ano (YYYY)"
// @Success 200 {array} models.Income
// @Failure 400,401,500 {string} string
// @Router /incomes/{userId} [get]
func ListIncomes(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["userId"]
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	if month < 1 || month > 12 || year == 0 {
		http.Error(w, "Parâmetros de mês ou ano inválidos", http.StatusBadRequest)
		return
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	rows, err := db.DB.Query(`
		SELECT id, user_id, descricao, valor, data_recebimento, categoria, observacoes, created_at
		FROM incomes
		WHERE user_id = $1 AND data_recebimento >= $2 AND data_recebimento < $3
		ORDER BY data_recebimento
	`, uid, start, end) // Consulta as receitas do usuário no mês e ano especificados

	if err != nil {
		http.Error(w, "Erro ao consultar receitas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close() // Fecha as linhas após a leitura

	var result []models.Income
	for rows.Next() {
		var inc models.Income
		if err := rows.Scan(&inc.ID, &inc.UserID, &inc.Descricao, &inc.Valor, &inc.DataRecebimento, &inc.Categoria, &inc.Observacoes, &inc.CreatedAt); err != nil {
			http.Error(w, "Erro ao ler receita: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, inc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetIncomeByID busca uma receita específica pelo ID
//
// @Summary Buscar receita por ID
// @Tags Incomes
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da receita"
// @Success 200 {object} models.Income
// @Failure 400,401,404,500 {string} string
// @Router /incomes/{userId}/{id} [get]
func GetIncomeByID(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	row := db.DB.QueryRow(`
		SELECT id, user_id, descricao, valor, data_recebimento, categoria, observacoes, created_at
		FROM incomes
		WHERE user_id=$1 AND id=$2
	`, p["userId"], p["id"])

	var inc models.Income
	if err := row.Scan(&inc.ID, &inc.UserID, &inc.Descricao, &inc.Valor, &inc.DataRecebimento, &inc.Categoria, &inc.Observacoes, &inc.CreatedAt); err != nil {
		http.Error(w, "Receita não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inc)
}

// Atualizar receita
//
// @Summary Atualizar receita
// @Tags Incomes
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da receita"
// @Param income body models.Income true "Dados da receita"
// @Success 200 {object} models.Income
// @Failure 400,401,500 {string} string
// @Router /incomes/{userId}/{id} [put]
func UpdateIncome(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	userID := p["userId"]
	incomeID := p["id"]

	var in struct {
		Descricao       string    `json:"descricao"`
		Valor           float64   `json:"valor"`
		DataRecebimento time.Time `json:"data_recebimento"`
		Categoria       string    `json:"categoria"`
		Observacoes     *string   `json:"observacoes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if in.Valor <= 0 {
		http.Error(w, "Valor deve ser positivo", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE incomes
		SET descricao=$1, valor=$2, data_recebimento=$3, categoria=$4, observacoes=$5
		WHERE user_id=$6 AND id=$7
		RETURNING id, user_id, descricao, valor, data_recebimento, categoria, observacoes, created_at;
	`

	var out models.Income
	err := db.DB.QueryRow(query,
		in.Descricao, in.Valor, in.DataRecebimento, in.Categoria, in.Observacoes,
		userID, incomeID).
		Scan(&out.ID, &out.UserID, &out.Descricao, &out.Valor,
			&out.DataRecebimento, &out.Categoria, &out.Observacoes, &out.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Receita não encontrada", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Erro ao atualizar receita: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// DeleteIncome exclui uma receita de um usuário pelo ID
//
// @Summary Excluir receita
// @Tags Incomes
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da receita"
// @Success 200 {object} map[string]string
// @Failure 400,401,404,500 {string} string
// @Router /incomes/{userId}/{id} [delete]
func DeleteIncome(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	res, err := db.DB.Exec(`
		DELETE FROM incomes
		WHERE user_id = $1 AND id = $2
	`, p["userId"], p["id"]) // Deleta a receita do usuário

	if err != nil {
		http.Error(w, "Erro ao deletar receita: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "Receita não encontrada", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Receita excluída com sucesso"})
}
