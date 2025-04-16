package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type ErrorCode string

const (
	ErrInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrMissingField  ErrorCode = "MISSING_FIELD"
	ErrInvalidFormat ErrorCode = "INVALID_FORMAT"
	ErrInvalidValue  ErrorCode = "INVALID_VALUE"
)

// ValidationError represents a validation error with additional context
type ValidationError struct {
	Code    ErrorCode
	Message string
	Field   string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	errorMap := map[string]interface{}{
		"message": e.Message,
		"code":    string(e.Code),
	}

	if e.Field != "" {
		errorMap["field"] = e.Field
	}

	if e.Value != nil {
		errorMap["value"] = e.Value
	}

	jsonBytes, err := json.Marshal(errorMap)
	if err != nil {
		return fmt.Sprintf("Error code: %s, Message: %s, Field: %s", e.Code, e.Message, e.Field)
	}

	return string(jsonBytes)
}

// NewValidationError creates a new validation error with the specified code and message
func NewValidationError(code ErrorCode, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Message: message,
	}
}

// WithField adds field information to the error
func (e *ValidationError) WithField(field string) *ValidationError {
	e.Field = field
	return e
}

// WithValue adds the invalid value to the error
func (e *ValidationError) WithValue(value interface{}) *ValidationError {
	e.Value = value
	return e
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// GetValidationError returns the ValidationError if the error is one, nil otherwise
func GetValidationError(err error) *ValidationError {
	if ve, ok := err.(*ValidationError); ok {
		return ve
	}
	return nil
}

// Common validation error constructors
func NewInvalidInputError(message string) *ValidationError {
	return NewValidationError(ErrInvalidInput, message)
}

func NewMissingFieldError(field string) *ValidationError {
	return NewValidationError(ErrMissingField, "missing required field").WithField(field)
}

func NewInvalidFormatError(field string, value interface{}) *ValidationError {
	return NewValidationError(ErrInvalidFormat, "invalid format").WithField(field).WithValue(value)
}

func NewInvalidValueError(field string, value interface{}) *ValidationError {
	return NewValidationError(ErrInvalidValue, "invalid value").WithField(field).WithValue(value)
}

func IsValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func IsValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

func IsValidCurrencyFormat(amount string) bool {
	currencyPattern := regexp.MustCompile(`^\d+\.\d{2}$`)
	return currencyPattern.MatchString(amount)
}
