package routes

import (
	"finance/src/controllers"
	"finance/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	r.HandleFunc("/refresh", controllers.RefreshToken).Methods("POST")

	// protegido: wrap com o middleware JWTAuth
	secure := middlewares.JWTAuth

	r.Handle("/users/{id}", secure(http.HandlerFunc(controllers.GetUserById))).Methods("GET")

	r.Handle("/expenses", secure(http.HandlerFunc(controllers.CreateExpense))).Methods("POST")
	// GET /users/{userId}?month=10&year=2023
	r.Handle("/expenses/{userId}", secure(http.HandlerFunc(controllers.ListExpenses))).Methods("GET")
	// Rota para listar todas as despesas de um usuário
	r.Handle("/expenses/{userId}", secure(http.HandlerFunc(controllers.ListAllExpenses))).Methods("GET")
	r.Handle("/users/{userId}/expenses/{id}", secure(http.HandlerFunc(controllers.GetExpenseByID))).Methods("GET")
	r.Handle("/users/{userId}/expenses/{id}", secure(http.HandlerFunc(controllers.UpdateExpense))).Methods("PUT")
	r.Handle("/users/{userId}/expenses/{id}", secure(http.HandlerFunc(controllers.DeleteExpense))).Methods("DELETE")
	r.Handle("/users/{userId}/expenses/{id}/pay", secure(http.HandlerFunc(controllers.PayExpense))).Methods("PATCH")
	r.Handle("/users/{userId}/expenses/{id}/unpay", secure(http.HandlerFunc(controllers.UnpayExpense))).Methods("PATCH")

	// Rota para obter o resumo mensal
	r.Handle("/summary/{userId}", secure(http.HandlerFunc(controllers.GetMonthlySummary))).Methods("GET")

	// Rota categorias
	r.Handle("/categories", secure(http.HandlerFunc(controllers.CreateCategory))).Methods("POST")
	r.Handle("/categories/{userId}", secure(http.HandlerFunc(controllers.GetCategories))).Methods("GET")
	r.Handle("/categories/{id}", secure(http.HandlerFunc(controllers.DeleteCategory))).Methods("DELETE")

	return r
}
