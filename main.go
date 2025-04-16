package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt-processor/handlers"
	"receipt-processor/logging"
	"receipt-processor/models"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the logger
	const logDir = "logs"
	logging.InitializeLogger(logDir)
	logging.LogInfo("Starting Receipt Processor application", logging.LogParams{})

	// Create a new router
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(loggingMiddleware)

	// Create a receipt store
	store := models.NewReceiptStore()

	// Create a receipt handler
	receiptHandler := handlers.NewReceiptHandler(store)

	// Register routes
	router.HandleFunc("/receipts/process", receiptHandler.ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", receiptHandler.GetPoints).Methods("GET")

	// Start the server
	port := 8080
	serverAddr := fmt.Sprintf(":%d", port)
	logging.LogInfo("Server starting", logging.LogParams{"port": port})
	fmt.Printf("Server starting on port %d...\n", port)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logging.LogError("Server failed to start", logging.LogParams{"error": err.Error()})
		log.Fatal(err)
	}
}

// loggingMiddleware logs each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.LogInfo("Request received", logging.LogParams{
			"method": r.Method,
			"path":   r.URL.Path,
			"remote": r.RemoteAddr,
		})
		next.ServeHTTP(w, r)
	})
}
