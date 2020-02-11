package jwt

import (
	"cronspy/backend/pkg/util/exception"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/parnurzeal/gorequest"
	"github.com/pquerna/ffjson/ffjson"
)

// Service provides a Json-Web-Token authentication implementation
// through Auth service
type Service struct {
	AuthURL     string
	HTTPTimeout time.Duration
}

// New generates new JWT service necessery for auth middleware
func New(authServiceURL string, httpTimeout time.Duration) *Service {
	return &Service{
		AuthURL:     authServiceURL,
		HTTPTimeout: httpTimeout,
	}
}

// MiddleWareFunc makes JWT implement the Middleware interface.
// It extracts and validates token from the request
func (j *Service) MiddleWareFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			retErr := false

			token := c.Request().Header.Get("Authorization")
			if token == "" {
				retErr = true
			}

			parts := strings.SplitN(token, " ", 2)
			if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
				retErr = true
			}

			// if something happened, return 401
			if retErr {
				return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeUnauthorized, ""))
			}

			// validate token using auth-service
			claims, err := j.ValidateToken(parts[1], c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeUnauthorized, ""))
			}

			// if token is valid, save claims into context
			for k, v := range claims {
				c.Set(k, v)
			}

			return next(c)
		}
	}
}

// ValidateToken parses token from Authorization header
func (j *Service) ValidateToken(token string, c echo.Context) (claims map[string]interface{}, err error) {

	// invoke auth service for token validation
	url := fmt.Sprintf("%s/validate_token", j.AuthURL)
	httpRequest := gorequest.New().Post(url)

	payload := fmt.Sprintf(`{ "token" : "%s" }`, token)

	if response, _, errRequest := httpRequest.Send(payload).Timeout(j.HTTPTimeout).End(); errRequest == nil {

		decoder := ffjson.NewDecoder()
		if response.StatusCode != http.StatusOK {

			if response.StatusCode == http.StatusUnauthorized {
				err = errors.New("invalid token")
			} else {
				err = errors.New("error validating access token")
			}

		} else {

			// retrieve claims
			type tokenClaim struct {
				Claims     map[string]interface{} `json:"claims"`
				Expiration int                    `json:"expiration"`
			}

			c := new(tokenClaim)
			if errJSON := decoder.DecodeReader(response.Body, c); errJSON == nil {
				claims = c.Claims
			}
		}

	} else {
		err = exception.ErrInternal
	}

	return
}

// HealthCheck calls the health check endpoint in auth-service
func (j *Service) HealthCheck() (err error) {

	url := fmt.Sprintf("%s/health", j.AuthURL)
	httpRequest := gorequest.New().Get(url)

	if _, _, errRequest := httpRequest.Timeout(j.HTTPTimeout).End(); errRequest != nil {
		err = fmt.Errorf("error calling apc-auth-service health check: %s", errRequest)
	}

	return
}
