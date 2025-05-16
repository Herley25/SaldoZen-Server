package controllers

import (
	"encoding/json"
	"finance/src/db"
	"finance/src/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Criar uma nova categoria
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if cat.Name == "" || cat.UserID == uuid.Nil {
		http.Error(w, "Nome e ID do usuário são obrigatórios", http.StatusBadRequest)
		return
	}

	cat.ID = uuid.New()
	cat.CreatedAt = time.Now()

	_, err := db.DB.Exec(`
		INSERT INTO categories (id, user_id, name, created_at)
		VALUES ($1, $2, $3, $4)
	`, cat.ID, cat.UserID, cat.Name, cat.CreatedAt)

	if err != nil {
		http.Error(w, "Erro ao criar categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cat)
}

// Listar todas as categorias de um usuário
func GetCategories(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	rows, err := db.DB.Query(`
		SELECT id, user_id, name, created_at
		FROM categories
		WHERE user_id = $1
		ORDER BY name ASC
	`, userID)
	if err != nil {
		http.Error(w, "Erro ao buscar categorias: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.CreatedAt); err != nil {
			http.Error(w, "Erro ao ler categoria: "+err.Error(), http.StatusInternalServerError)
			return
		}
		categories = append(categories, cat)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}

// Excluir uma categoria
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	result, err := db.DB.Exec(`
		DELETE FROM categories
		WHERE id = $1
	`, id)
	if err != nil {
		http.Error(w, "Erro ao excluir categoria: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Nenhuma categoria encontrada com esse ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Categoria excluída com sucesso"))
}
