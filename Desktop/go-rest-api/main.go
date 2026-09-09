package main

import (
	"fmt"
	"net/http"

	"go-rest-api/config"
	"go-rest-api/handlers"
)

func main() {
	config.ConnectDB()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /expenses", handlers.CreateExpense)
	mux.HandleFunc("GET /expenses", handlers.GetExpenses)
	mux.HandleFunc("GET /expenses/summary", handlers.GetSummaryByCategory)
	mux.HandleFunc("GET /expenses/{id}", handlers.GetExpenseByID)
	mux.HandleFunc("PATCH /expenses/{id}", handlers.UpdateExpense)
	mux.HandleFunc("DELETE /expenses/{id}", handlers.DeleteExpense)


	fmt.Println("Server 8080 portunda işləyir...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server dayandı xəta:", err)
	}
}