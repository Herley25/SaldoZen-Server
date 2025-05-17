package controllers

import (
	"encoding/json"
	"finance/src/db"
	"finance/src/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Nome, e-mail e senha são obrigatórios", http.StatusBadRequest)
		return
	}

	// Gerar o hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao gerar hash da senha: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = uuid.New()
	user.PasswordHash = string(hashedPassword)
	user.CreatedAt = time.Now()

	_, err = db.DB.Exec(`
		INSERT INTO users (id, name, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt)

	if err != nil {
		http.Error(w, "Erro ao salvar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// limpar a senha do retorno
	user.Password = ""

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
