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
	CodeInternalServerError     = "internal_server_error"
	CodeUnauthorized            = "unauthorized"
	CodeUserExists              = "user_exists"
	CodeInvalidPassword         = "invalid_password"
	CodeUnknownUser             = "unknown_user"
	CodeMaxPasswordResetReached = "max_password_reset_reached"
	CodeNeedToWaitBeforeResend  = "wait_before_resend"
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
		CodeUnknownUser:                  "unknown user",
		CodeMaxPasswordResetReached:      "max number of password resets has been reached",
		CodeNeedToWaitBeforeResend:       "need to wait some time before resending the email",
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
