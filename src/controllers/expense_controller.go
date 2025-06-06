package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"finance/src/db"
	"finance/src/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Controller para criar uma nova despesa

// @Summary	Criar despesa
// @Tags		Expenses
// @Security	BearerAuth
// @Accept		json
// @Produce	json
// @Param		expense	body		models.Expense	true	"Despesa"
// @Success	201	{object}	models.Expense
// @Failure	400,401	{string}	string
// @Router		/expenses [post]
func CreateExpense(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if expense.Valor <= 0 {
		http.Error(w, "Valor deve ser maior que zero", http.StatusBadRequest)
		return
	}

	expense.ID = uuid.New()
	expense.CreatedAt = time.Now()
	if expense.Paga && expense.DataPagamento == nil {
		now := time.Now()
		expense.DataPagamento = &now
	}

	_, err := db.DB.Exec(`
		INSERT INTO expenses (id, user_id, descricao, valor, vencimento, paga, data_pagamento, categoria, observacoes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, expense.ID, expense.UserID, expense.Descricao, expense.Valor, expense.Vencimento, expense.Paga, expense.DataPagamento, expense.Categoria, expense.Observacoes, expense.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}

// ListExpenses retorna as despesas de um usuário filtradas por mês e ano
//
// @Summary Listar despesas
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param month query string false "Mês (1-12)"
// @Param year query string false "Ano (YYYY)"
// @Success 200 {array} models.Expense
// @Failure 400,401,500 {string} string
// @Router /expenses/{userId} [get]
func ListExpenses(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"]

	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	var filters []interface{}
	query := `
		SELECT id, user_id, descricao, valor, vencimento, paga, data_pagamento, categoria, observacoes, created_at
		FROM expenses
		WHERE user_id = $1
	`
	filters = append(filters, userId)

	if monthStr != "" && yearStr != "" {
		month, err1 := strconv.Atoi(monthStr)
		year, err2 := strconv.Atoi(yearStr)

		if err1 != nil || err2 != nil {
			http.Error(w, "Parâmetros de mês e ano inválidos", http.StatusBadRequest)
			return
		}

		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		query += " AND vencimento >= $2 AND vencimento < $3"
		filters = append(filters, startDate, endDate)
	}

	rows, err := db.DB.Query(query, filters...)
	if err != nil {
		http.Error(w, "Erro ao buscar despesas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var dataPagamento sql.NullTime
		var observacoes sql.NullString

		err := rows.Scan(&e.ID, &e.UserID, &e.Descricao, &e.Valor, &e.Vencimento, &e.Paga, &dataPagamento, &e.Categoria, &observacoes, &e.CreatedAt)
		if err != nil {
			http.Error(w, "Erro ao ler despesa: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if dataPagamento.Valid {
			e.DataPagamento = &dataPagamento.Time
		}
		if observacoes.Valid {
			e.Observacoes = &observacoes.String
		}

		expenses = append(expenses, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// ListAllExpenses retorna todas as despesas de um usuário
//
// @Summary Listar todas as despesas
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Success 200 {array} models.Expense
// @Failure 400,401,500 {string} string
// @Router /expenses/{userId} [get]
func ListAllExpenses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]

	rows, err := db.DB.Query(`
		SELECT id, user_id, descricao, valor, vencimento, paga, data_pagamento, categoria, observacoes, created_at
		FROM expenses
		WHERE user_id = $1
	`, userId)
	if err != nil {
		http.Error(w, "Erro ao buscar despesas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.UserID, &e.Descricao, &e.Valor, &e.Vencimento, &e.Paga, &e.DataPagamento, &e.Categoria, &e.Observacoes, &e.CreatedAt)

		if err != nil {
			http.Error(w, "Erro ao ler despesa: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.Expense{
			ID:            e.ID,
			UserID:        e.UserID,
			Descricao:     e.Descricao,
			Valor:         e.Valor,
			Vencimento:    e.Vencimento,
			Paga:          e.Paga,
			DataPagamento: e.DataPagamento,
			Categoria:     e.Categoria,
			Observacoes:   e.Observacoes,
			CreatedAt:     e.CreatedAt,
		}

		expenses = append(expenses, response)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

// GetExpenseByID busca uma despesa específica pelo ID
//
// @Summary Buscar despesa por ID
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param expenseId path string true "ID da despesa"
// @Success 200 {object} models.Expense
// @Failure 400,401,404,500 {string} string
// @Router /users/{userId}/expenses/{id} [get]
func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	expenseId := params["expenseId"]

	row := db.DB.QueryRow(`
		SELECT id, user_id, descricao, valor, vencimento, paga, data_pagamento, categoria, observacoes, created_at
		FROM expenses
		WHERE user_id = $1 AND id = $2
	`, userId, expenseId)

	var e models.Expense
	err := row.Scan(
		&e.ID,
		&e.UserID,
		&e.Descricao,
		&e.Valor,
		&e.Vencimento,
		&e.Paga,
		&e.DataPagamento,
		&e.Categoria,
		&e.Observacoes,
		&e.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Despesa não encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Erro ao buscar despesa: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := models.Expense{
		ID:            e.ID,
		UserID:        e.UserID,
		Descricao:     e.Descricao,
		Valor:         e.Valor,
		Vencimento:    e.Vencimento,
		Paga:          e.Paga,
		DataPagamento: e.DataPagamento,
		Categoria:     e.Categoria,
		Observacoes:   e.Observacoes,
		CreatedAt:     e.CreatedAt,
		Status:        e.StatusHoje(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateExpense atualiza uma despesa existente
//
// @Summary Atualizar despesa
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da despesa"
// @Param expense body models.Expense true "Dados da despesa"
// @Success 204 {string} string "Despesa atualizada com sucesso"
// @Failure 400,401,404,500 {string} string
// @Router /users/{userId}/expenses/{id} [put]
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	expenseId := params["id"]

	var update models.Expense
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if update.Valor <= 0 {
		http.Error(w, "Valor deve ser maior que zero", http.StatusBadRequest)
		return
	}

	// Se a despesa foi paga e sem data, define a data de pagamento como a data atual
	if update.Paga && update.DataPagamento == nil {
		now := time.Now()
		update.DataPagamento = &now
	}

	// Se desmarcar a despesa como paga, remove a data de pagamento
	if !update.Paga {
		update.DataPagamento = nil
	}

	result, err := db.DB.Exec(`
		UPDATE expenses
		SET descricao = $1, valor = $2, vencimento = $3, paga = $4, data_pagamento = $5, categoria = $6, observacoes = $7
		WHERE user_id = $8 AND id = $9
	`, update.Descricao, update.Valor, update.Vencimento, update.Paga, update.DataPagamento, update.Categoria, update.Observacoes, userId, expenseId)

	if err != nil {
		http.Error(w, "Erro ao atualizar despesa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Despesa não encontrada ou não pertence ao usuário", http.StatusNotFound)
		return
	}

	// Retorna a despesa atualizada
	w.WriteHeader(http.StatusNoContent)
}

// Excluir uma despesa
type Message struct {
	Message string `json:"message"`
}

// DeleteExpense exclui uma despesa de um usuário pelo ID
//
// @Summary Excluir despesa
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da despesa"
// @Success 200 {object} string "Despesa excluída com sucesso"
// @Failure 400,401,404,500 {string} string
// @Router /users/{userId}/expenses/{id} [delete]
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	expenseId := params["id"]

	result, err := db.DB.Exec(`
		DELETE FROM expenses
		WHERE user_id = $1 AND id = $2
	`, userId, expenseId)

	if err != nil {
		http.Error(w, "Erro ao excluir despesa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Despesa não encontrada ou não pertence ao usuário", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Message{Message: "Despesa excluída com sucesso"})
}

// PayExpense marca uma despesa como paga
//
// @Summary Marcar despesa como paga
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da despesa"
// @Success 200 {object} Message "Despesa marcada como paga com sucesso"
// @Failure 400,401,404,500 {string} string
// @Router /users/{userId}/expenses/{id}/pay [patch]
func PayExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	expenseId := params["id"]

	now := time.Now()

	// Atualiza a despesa para marcada como paga
	result, err := db.DB.Exec(`
		UPDATE expenses
		SET paga = true, data_pagamento = $1
		WHERE user_id = $2 AND id = $3
	`, now, userId, expenseId)

	if err != nil {
		http.Error(w, "Erro ao marcar despesa como paga: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Despesa não encontrada ou não pertence ao usuário", http.StatusNotFound)
		return
	}

	// Retorna a despesa atualizada
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Message{Message: "Despesa marcada como paga com sucesso"})
}

// UnpayExpense marca uma despesa como não paga
//
// @Summary Marcar despesa como não paga
// @Tags Expenses
// @Security BearerAuth
// @Param userId path string true "ID do usuário"
// @Param id path string true "ID da despesa"
// @Success 200 {object} Message "Despesa marcada como não paga com sucesso"
// @Failure 400,401,404,500 {string} string
// @Router /users/{userId}/expenses/{id}/unpay [patch]
func UnpayExpense(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	expenseId := params["id"]

	// Atualiza a despesa para marcada como não paga
	result, err := db.DB.Exec(`
		UPDATE expenses
		SET paga = false, data_pagamento = null
		WHERE user_id = $1 AND id = $2
	`, userId, expenseId)

	if err != nil {
		http.Error(w, "Erro ao marcar despesa como não paga: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Despesa não encontrada ou não pertence ao usuário", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Message{Message: "Despesa marcada como não paga com sucesso"})
}
