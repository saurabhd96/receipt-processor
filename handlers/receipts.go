package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"receipt-processor/logging"
	"receipt-processor/models"

	"github.com/gorilla/mux"
)

// ReceiptHandler handles receipt-related requests
type ReceiptHandler struct {
	Store *models.ReceiptStore
}

// NewReceiptHandler creates a new receipt handler
func NewReceiptHandler(store *models.ReceiptStore) *ReceiptHandler {
	return &ReceiptHandler{
		Store: store,
	}
}

// ProcessReceipt processes a receipt and returns an ID
func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	logging.LogInfo("Processing receipt request", logging.LogParams{
		"method": r.Method,
		"path":   r.URL.Path,
	})

	// Read and parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.LogError("Failed to read request body", logging.LogParams{
			"error": err.Error(),
		})
		http.Error(w, "Error reading request body. Please verify input.", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Log the request body for debugging
	logging.LogInfo("Receipt data received", logging.LogParams{
		"data": string(body),
	})

	// Parse JSON
	var receipt models.Receipt
	err = json.Unmarshal(body, &receipt)
	if err != nil {
		logging.LogError("Failed to parse JSON", logging.LogParams{
			"error": err.Error(),
		})
		http.Error(w, "Invalid JSON format. Please verify input.", http.StatusBadRequest)
		return
	}

	// Validate receipt
	if err := validateReceipt(receipt); err != nil {
		logging.LogWarn("Invalid receipt data", logging.LogParams{
			"error": string(err.Error()),
		})
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate points for the receipt
	points := models.CalculatePoints(receipt)
	receipt.Points = points

	logging.LogInfo("Points calculated", logging.LogParams{
		"points": points,
	})

	// Store the receipt with calculated points
	id := h.Store.AddReceipt(receipt)

	logging.LogInfo("Receipt processed", logging.LogParams{
		"id":     id,
		"points": points,
	})

	// Prepare and send response
	response := models.ProcessResponse{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.LogError("Failed to encode response", logging.LogParams{
			"error": err.Error(),
		})
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
}

// GetPoints returns the points for a receipt
func (h *ReceiptHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	params := mux.Vars(r)
	id := params["id"]

	logging.LogInfo("Getting points for receipt", logging.LogParams{
		"id": id,
	})

	// Get receipt from store
	receipt, found := h.Store.GetReceipt(id)
	if !found {
		logging.LogWarn("Receipt not found", logging.LogParams{
			"id": id,
		})
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Prepare response
	response := models.PointsResponse{
		Points: receipt.Points,
	}

	logging.LogInfo("Returning points", logging.LogParams{
		"id":     id,
		"points": receipt.Points,
	})

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.LogError("Failed to encode response", logging.LogParams{
			"error": err.Error(),
		})
		http.Error(w, "Failed to generate response", http.StatusInternalServerError)
		return
	}
}

// validateReceipt validates the receipt data
func validateReceipt(receipt models.Receipt) error {
	// Basic validation with proper error types
	if receipt.Retailer == "" {
		return models.NewMissingFieldError("retailer")
	}
	if receipt.PurchaseDate == "" {
		return models.NewMissingFieldError("purchaseDate")
	}
	if receipt.PurchaseTime == "" {
		return models.NewMissingFieldError("purchaseTime")
	}
	if receipt.Total == "" {
		return models.NewMissingFieldError("total")
	}
	if len(receipt.Items) == 0 {
		return models.NewMissingFieldError("items")
	}

	// Validate purchase date format (YYYY-MM-DD)
	if !models.IsValidDateFormat(receipt.PurchaseDate) {
		return models.NewInvalidFormatError("purchaseDate", receipt.PurchaseDate)
	}

	// Validate purchase time format (HH:MM)
	if !models.IsValidTimeFormat(receipt.PurchaseTime) {
		return models.NewInvalidFormatError("purchaseTime", receipt.PurchaseTime)
	}

	// Validate total is a valid number
	if !models.IsValidCurrencyFormat(receipt.Total) {
		return models.NewInvalidFormatError("total", receipt.Total)
	}

	// Validate each item
	for i, item := range receipt.Items {
		if item.ShortDescription == "" {
			return models.NewMissingFieldError("items[" + string(rune(i)) + "].shortDescription")
		}
		if item.Price == "" {
			return models.NewMissingFieldError("items[" + string(rune(i)) + "].price")
		}
		if !models.IsValidCurrencyFormat(item.Price) {
			return models.NewInvalidFormatError("items["+string(rune(i))+"].price", item.Price)
		}
	}

	return nil
}
