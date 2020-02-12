package transport

import (
	"cronspy/backend/pkg/api/user"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"
	"time"

	"github.com/astropay/go-tools/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	jwtSigningKey = "10fa4f27-6a69-45c1-9a88-dfcecdbdc3d8"
)

var (
	// IsUserLoggedIn is a middleware to restrict URL to logged user
	IsUserLoggedIn = middleware.JWTWithConfig(getJWTConfig())
	// used jwt signing method
	jwtSigningMethod = jwt.SigningMethodHS512
)

// HTTP represents auth http service
type HTTP struct {
	svc user.Service
}

// NewHTTP creates new http service
func NewHTTP(svc user.Service, e *echo.Echo) {
	h := HTTP{svc: svc}

	user := e.Group("/user")
	user.POST("/register", h.userRegisterHandler)
	user.POST("/login", h.userLoginHandler)
	user.PUT("/changePassword", h.userChangePasswordHandler, IsUserLoggedIn)
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
	if err := h.validateEmailInput(user.Email); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(err.Error(), ""))
	}
	if err := h.validatePasswordInput(user.Password); err != nil {
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
	if err := h.validateEmailInput(cred.Username); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
	}
	if err := h.validatePasswordInput(cred.Password); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
	}

	// run login
	user, err := h.svc.Login(cred.Username, cred.Password)
	if err != nil {
		return err
	}

	// generate JWT
	token := h.buildJWTToken(user.ID, user.Email, user.Name, user.AccountType, h.svc.GetJWTExpiration())
	t, errSign := token.SignedString([]byte(jwtSigningKey))
	if errSign != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errSign.Error()))
	}

	resp := make(map[string]interface{})
	resp["user"] = user
	resp["access_token"] = t

	return c.JSON(http.StatusOK, resp)
}

//
// CHANGE PASSWORD
//
func (h *HTTP) userChangePasswordHandler(c echo.Context) error {

	type credentials struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	cred := new(credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}

	// if email or password does not meet basic criteria, return generic error
	if err := h.validatePasswordInput(cred.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.ErrInvalidPasswordFormat.Error(), ""))
	}

	// get user id
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idUser, ok := claims["id"].(float64)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
	}

	// run change password
	err := h.svc.ChangePassword(int(idUser), cred.OldPassword, cred.NewPassword)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *HTTP) validateEmailInput(email string) (err error) {
	// check email format
	if !common.IsEmailAddress(email) {
		return exception.ErrInvalidEmailAddress
	}

	return
}

func (h *HTTP) validatePasswordInput(password string) (err error) {
	// check password rules
	if len(password) < 8 {
		return exception.ErrInvalidPasswordFormat
	}

	return
}

// build JWT with the indicated parameters
func (h *HTTP) buildJWTToken(userID int, email, name, accountType string, tokenExpiration int) *jwt.Token {
	token := jwt.New(jwtSigningMethod)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["email"] = email
	claims["name"] = name
	claims["account_type"] = accountType
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenExpiration)).Unix()

	return token
}

func getJWTConfig() (jwtCfg middleware.JWTConfig) {
	jwtCfg.SigningMethod = jwtSigningMethod.Name
	jwtCfg.SigningKey = []byte(jwtSigningKey)
	return
}
