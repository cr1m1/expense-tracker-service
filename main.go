package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cr1m1/expense-tracker-service/database"
	"github.com/cr1m1/expense-tracker-service/handlers"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	db, err := database.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	handler := handlers.NewExpenseHandler(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateExpense(w, r)
		} else if r.Method == http.MethodGet {
			if id := r.URL.Query().Get("id"); id != "" {
				handler.GetExpense(w, r)
			} else {
				handler.ListExpenses(w, r)
			}
		} else if r.Method == http.MethodPut {
			handler.UpdateExpense(w, r)
		} else if r.Method == http.MethodDelete {
			handler.DeleteExpense(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8002"
	fmt.Printf("Starting expense tracker service on %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
