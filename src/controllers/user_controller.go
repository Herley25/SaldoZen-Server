package controllers

import (
	"encoding/json"
	"finance/src/db"
	"finance/src/models"
	"finance/src/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
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

// LoginUser faz o login de um usuário
//
//	@Summary	Login
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		credentials	body		object{email=string,password=string}	true	"Credenciais"
//	@Success	200		{object}	map[string]string
//	@Failure	400,401	{string}	string	"erro"
//	@Router		/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`, creds.Email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Erro ao buscar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verifica a senha
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)) != nil {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// Cria um token JWT
	secret := []byte(os.Getenv("JWT_SECRET"))
	expHours := 24
	if h := os.Getenv("JWT_EXP_HOURS"); h != "" {
		if v, err := strconv.Atoi(h); err == nil {
			expHours = v
		}
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Duration(expHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString(secret)

	access, _ := utils.GenerateAccess(user.ID.String())   // Gera o access token
	refresh, _ := utils.GenerateRefresh(user.ID.String()) // Gera o refresh token

	http.SetCookie(w, &http.Cookie{ // Define o cookie de refresh token
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]string{ // Retorna o token JWT e o access token
		"token":        signed,
		"access_token": access,
	})
}

// GetUser retorna os dados do usuário autenticado
func GetUserById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var user models.User
	err := db.DB.QueryRow(`
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Erro ao buscar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// RefreshToken atualiza o token JWT
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token não encontrado", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	secret := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { // Verifica o método de assinatura
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token inválido: "+err.Error(), http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	newAccessToken, _ := utils.GenerateAccess(userID) // Gera o novo access token

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
