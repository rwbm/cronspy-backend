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
	user.POST("/login", h.userLoginHandler)
}

//
// USER REGISTRATION
//
func (h *HTTP) userRegisterHandler(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	// validate input
	if err := h.validateUserRegistrationInput(user.Email, user.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(err.Error(), ""))
	}

	// register user
	if err := h.svc.RegisterUser(c, user); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

//
// LOGIN
//
func (h *HTTP) userLoginHandler(c echo.Context) error {

	type credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	cred := new(credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}

	// if email or password does not meet basic criteria, return generic error
	if err := h.validateUserRegistrationInput(cred.Username, cred.Password); err != nil {
		err = echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
	}

	// run login
	user, err := h.svc.Login(cred.Username, cred.Password)
	if err != nil {
		return err
	}

	//
	// TODO: ademas debemos retornar el token JWT
	//

	return c.JSON(http.StatusOK, user)
}

func (h *HTTP) validateUserRegistrationInput(email, password string) (err error) {

	// check email format
	if !common.IsEmailAddress(email) {
		return exception.ErrInvalidEmailAddress
	}

	// check password rules
	if len(password) < 8 {
		return exception.ErrInvalidPasswordFormat
	}

	return
}
