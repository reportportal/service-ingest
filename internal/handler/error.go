package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`              // user-level status message
	AppCode    int64  `json:"errorCode,omitempty"` // application-specific error code
	ErrorText  string `json:"message,omitempty"`   // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		AppCode:        4001,
		ErrorText:      "Incorrect Request. " + parseValidationError(err),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var InternalServerError = &ErrResponse{HTTPStatusCode: 500, StatusText: "Internal server error."}
var NotImplementedError = &ErrResponse{HTTPStatusCode: 501, StatusText: "Not implemented."}

func RespondNotImplemented(w http.ResponseWriter, r *http.Request) {
	_ = render.Render(w, r, NotImplementedError)
}

func parseValidationError(err error) string {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, e := range validationErrs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", field)
			case "uuid":
				return fmt.Sprintf("%s must be a valid UUID", field)
			case "oneof":
				return fmt.Sprintf("%s must be one of: %s", field, e.Param())
			case "max":
				return fmt.Sprintf("%s must be less or equal than %s", field, e.Param())
			case "min":
				return fmt.Sprintf("%s must be greater or equal than %s", field, e.Param())
			default:
				return fmt.Sprintf("%s validation failed on %s", field, e.Tag())
			}
		}
	}
	return err.Error()
}
