package models

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Item represents an item on a receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt represents a receipt with purchasing information
type Receipt struct {
	ID           string `json:"id,omitempty"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
	Points       int    `json:"points,omitempty"`
}

// ProcessResponse is the response from the process endpoint
type ProcessResponse struct {
	ID string `json:"id"`
}

// PointsResponse is the response from the points endpoint
type PointsResponse struct {
	Points int `json:"points"`
}

// ReceiptStore is an in-memory store for receipts
type ReceiptStore struct {
	Receipts map[string]Receipt
}

// NewReceiptStore creates a new receipt store
func NewReceiptStore() *ReceiptStore {
	return &ReceiptStore{
		Receipts: make(map[string]Receipt),
	}
}

// AddReceipt adds a receipt to the store and returns the ID
func (rs *ReceiptStore) AddReceipt(receipt Receipt) string {
	id := uuid.New().String()
	receipt.ID = id
	rs.Receipts[id] = receipt
	return id
}

// GetReceipt gets a receipt from the store by ID
func (rs *ReceiptStore) GetReceipt(id string) (Receipt, bool) {
	receipt, ok := rs.Receipts[id]
	return receipt, ok
}

// CalculatePoints calculates the points for a receipt
func CalculatePoints(receipt Receipt) int {
	var points int

	// Rule 1: One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	if isRoundDollarAmount(receipt.Total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if isMultipleOf25Cents(receipt.Total) {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 && len(trimmedDesc) > 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				pointsForItem := int(math.Ceil(price * 0.2))
				points += pointsForItem
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd
	if isDayOdd(receipt.PurchaseDate) {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm
	if isTimeBetween2And4PM(receipt.PurchaseTime) {
		points += 10
	}

	return points
}

// countAlphanumeric counts the number of alphanumeric characters in a string
func countAlphanumeric(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// isRoundDollarAmount checks if the total is a round dollar amount
func isRoundDollarAmount(total string) bool {
	re := regexp.MustCompile(`^\d+\.00`)
	return re.MatchString(total)
}

// isMultipleOf25Cents checks if the total is a multiple of 0.25
func isMultipleOf25Cents(total string) bool {
	val, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return false
	}

	// Convert to cents and check if it's a multiple of 25
	cents := int(val * 100)
	return cents%25 == 0
}

// isDayOdd checks if the day in the purchase date is odd
func isDayOdd(purchaseDate string) bool {
	t, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		return false
	}

	day := t.Day()
	return day%2 == 1
}

// isTimeBetween2And4PM checks if the time is between 2:00PM and 4:00PM
func isTimeBetween2And4PM(purchaseTime string) bool {
	t, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return false
	}

	hour := t.Hour()
	minute := t.Minute()

	timeInMinutes := hour*60 + minute

	// 2:00PM = 14:00 = 14*60 = 840 minutes
	// 4:00PM = 16:00 = 16*60 = 960 minutes
	return timeInMinutes > 840 && timeInMinutes < 960
}
