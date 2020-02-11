package exception

import "errors"

// Internal error defitions
var (
	ErrRecordNotFound        = errors.New("record not found")
	ErrInvalidEmailAddress   = errors.New("invalid_email")
	ErrInvalidPasswordFormat = errors.New("invalid_password_format")
)

// Standar error codes
const (
	CodeInternalServerError = "internal_server_error"
	CodeUnauthorized        = "unauthorized"
	CodeUserExists          = "user_exists"
	CodeInvalidPassword     = "invalid_password"
)

var (
	// map with text messages associated
	messages = map[string]string{
		CodeInternalServerError:          "internal server error ocurred",
		CodeUnauthorized:                 "unauthorized request",
		CodeUserExists:                   "user already exists",
		ErrInvalidEmailAddress.Error():   "invalid email address",
		ErrInvalidPasswordFormat.Error(): "invalid password format",
		CodeInvalidPassword:              "invalid password",
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
