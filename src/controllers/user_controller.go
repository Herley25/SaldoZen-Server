package controllers

import (
	"encoding/json"
	"finance/src/db"
	"finance/src/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.PasswordHash == "" {
		http.Error(w, "Nome, e-mail e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err := db.DB.Exec(`
		INSERT INTO users (id, nome, email, senha, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt)

	if err != nil {
		http.Error(w, "Erro ao salvar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
