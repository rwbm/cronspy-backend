package transport

import (
	"cronspy/backend/pkg/api/user"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"

	"github.com/astropay/go-tools/common"
	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new http service
func NewHTTP(svc user.Service, e *echo.Echo) {
	h := HTTP{svc: svc}

	user := e.Group("/user")
	user.POST("/create", h.userRegisterHandler)
}

func (h *HTTP) userRegisterHandler(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	// validate input
	if err := h.validateUserRegistrationInput(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(err.Error(), ""))
	}

	// register user
	if err := h.svc.RegisterUser(c, user); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *HTTP) validateUserRegistrationInput(user *model.User) (err error) {

	// check email format
	if !common.IsEmailAddress(user.Email) {
		return exception.ErrInvalidEmailAddress
	}

	// check
	if len(user.Password) < 8 {
		return exception.ErrInvalidPasswordFormat
	}

	return
}
