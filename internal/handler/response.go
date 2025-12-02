package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// respondError writes error response
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, ErrorResponse{
		Error:   http.StatusText(code),
		Message: message,
	})
}

// respondValidationError writes validation error response
func respondValidationError(w http.ResponseWriter, err error) {
	details := make(map[string]string)

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				details[field] = "This field is required"
			case "uuid":
				details[field] = "Must be a valid UUID"
			case "oneof":
				details[field] = "Must be one of: " + e.Param()
			case "max":
				details[field] = "Maximum length is " + e.Param()
			case "min":
				details[field] = "Minimum length is " + e.Param()
			default:
				details[field] = "Validation failed on " + e.Tag()
			}
		}
	}

	respondJSON(w, http.StatusBadRequest, ErrorResponse{
		Error:   "Validation Error",
		Message: "Request validation failed",
		Details: details,
	})
}

// respondNotImplemented writes a 501 Not Implemented response
func respondNotImplemented(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
