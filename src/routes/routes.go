package routes

import (
	"net/http"

	"finance/src/controllers"

	"github.com/gorilla/mux"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")

	r.HandleFunc("/expenses", controllers.CreateExpense).Methods("POST")
	// GET /users/{userId}?month=10&year=2023
	r.HandleFunc("/expenses/{userId}", controllers.ListExpenses).Methods("GET")
	// Rota para listar todas as despesas de um usuário
	r.HandleFunc("/expenses/{userId}", controllers.ListAllExpenses).Methods("GET")
	r.HandleFunc("/users/{userId}/expenses/{id}", controllers.GetExpenseByID).Methods("GET")
	r.HandleFunc("/users/{userId}/expenses/{id}", controllers.UpdateExpense).Methods("PUT")
	r.HandleFunc("/users/{userId}/expenses/{id}", controllers.DeleteExpense).Methods("DELETE")
	r.HandleFunc("/users/{userId}/expenses/{id}/pay", controllers.PayExpense).Methods("PATCH")
	r.HandleFunc("/users/{userId}/expenses/{id}/unpay", controllers.UnpayExpense).Methods("PATCH")

	// Rota para obter o resumo mensal
	r.HandleFunc("/summary/{userId}", controllers.GetMonthlySummary).Methods("GET")

	// Rota categorias
	r.HandleFunc("/categories", controllers.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{userId}", controllers.GetCategories).Methods("GET")
	r.HandleFunc("/categories/{id}", controllers.DeleteCategory).Methods("DELETE")

	return r
}
