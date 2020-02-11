package exception

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Standar error codes
const (
	CodeInternalServerError = "internal_server_error"
	CodeUnauthorized        = "unauthorized"
)

var (
	// ErrInternal is used for generic internal server errors
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "something went wrong")

	// ErrNoResults is used to indicate no results from a DB query
	ErrNoResults = echo.NewHTTPError(http.StatusNotFound, "no results found")

	// ErrBadRequest (400) is returned for bad request (validation)
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest, "invalid or missing parameters")

	// ErrUnauthorized (401) is returned when user is not authorized
	ErrUnauthorized = echo.ErrUnauthorized
)

var (
	// map with text messages associated
	messages = map[string]string{
		CodeInternalServerError: "internal server error ocurred",
		CodeUnauthorized:        "unauthorized request",
	}
)

// GetErrorMap returns a map with the provided error code and associated message;
// it's useful for building HTTP error responses
func GetErrorMap(code, msg string) (m map[string]string) {
	if code != "" || msg != "" {
		m = make(map[string]string)

		if code != "" {
			m["code"] = code
		} else {
			m["code"] = CodeInternalServerError
		}

		if msg != "" {
			m["message"] = msg
		} else if code != "" {
			if message, ok := messages[code]; ok {
				m["message"] = message
			}
		}
	}
	return
}
