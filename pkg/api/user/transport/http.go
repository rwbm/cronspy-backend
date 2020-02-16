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

var (
	// IsUserLoggedIn is a middleware to restrict URL to logged user
	IsUserLoggedIn echo.MiddlewareFunc
)

// HTTP represents auth http service
type HTTP struct {
	svc              user.Service
	jwtSigningKey    string
	jwtSigningMethod *jwt.SigningMethodHMAC
}

// NewHTTP creates new http service to handle request to /user
func NewHTTP(svc user.Service, jwtSigningKey string, jwtSigningMethod *jwt.SigningMethodHMAC, e *echo.Echo) {
	h := HTTP{
		svc:              svc,
		jwtSigningKey:    jwtSigningKey,
		jwtSigningMethod: jwtSigningMethod,
	}

	// define logged user check function
	IsUserLoggedIn = middleware.JWTWithConfig(h.getJWTConfig())

	user := e.Group("/user")

	// --- Auth NOT required ---
	user.POST("/register", h.userRegisterHandler)
	user.POST("/login", h.userLoginHandler)

	user.POST("/passwordReset", h.userPasswordResetRequestHandler)
	user.GET("/passwordReset/validate", h.userPasswordResetValidateHandler)
	user.POST("/passwordReset/change", h.userPasswordResetChangeHandler)

	// --- Auth required ---
	user.PUT("/changePassword", h.userChangePasswordHandler, IsUserLoggedIn)
}

//
// --- USER REGISTRATION ---
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
// --- LOGIN ---
//
func (h *HTTP) userLoginHandler(c echo.Context) error {

	type credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	payload := new(credentials)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// if email or password does not meet basic criteria, return generic error
	if err := h.validateEmailInput(payload.Username); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
	}
	if err := h.validatePasswordInput(payload.Password); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
	}

	// run login
	user, err := h.svc.Login(payload.Username, payload.Password)
	if err != nil {
		return err
	}

	// generate JWT
	token := h.buildJWTToken(user.ID, user.Email, user.Name, user.AccountType, h.svc.GetJWTExpiration())
	t, errSign := token.SignedString([]byte(h.jwtSigningKey))
	if errSign != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errSign.Error()))
	}

	resp := make(map[string]interface{})
	resp["user"] = user
	resp["access_token"] = t

	return c.JSON(http.StatusOK, resp)
}

//
// --- PASSWORD RESET: REQUEST ---
//
func (h *HTTP) userPasswordResetRequestHandler(c echo.Context) error {

	type requestResponse struct {
		Email string `json:"email,omitempty"`
		ID    string `json:"id,omitempty"`
	}

	payload := new(requestResponse)
	if err := c.Bind(payload); err != nil {
		return err
	}

	resetID, err := h.svc.ResetPassword(payload.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, requestResponse{ID: resetID})
}

//
// --- PASSWORD RESET: VALIDATE ---
//
func (h *HTTP) userPasswordResetValidateHandler(c echo.Context) error {

	token := c.QueryParam("token")
	if token == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	err := h.svc.ValidateResetPassword(token)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

//
// --- PASSWORD RESET: CHANGE ---
//
func (h *HTTP) userPasswordResetChangeHandler(c echo.Context) error {

	type credentials struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	payload := new(credentials)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// if email or password does not meet basic criteria, return error
	if err := h.validatePasswordInput(payload.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.ErrInvalidPasswordFormat.Error(), ""))
	}

	// change password
	if err := h.svc.ChangePasswordWithReset(payload.Token, payload.NewPassword); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

//
// --- CHANGE PASSWORD ---
//
func (h *HTTP) userChangePasswordHandler(c echo.Context) error {

	type credentials struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	payload := new(credentials)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// if email or password does not meet basic criteria, return error
	if err := h.validatePasswordInput(payload.NewPassword); err != nil {
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
	err := h.svc.ChangePassword(int(idUser), payload.OldPassword, payload.NewPassword)
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
	token := jwt.New(h.jwtSigningMethod)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["email"] = email
	claims["name"] = name
	claims["account_type"] = accountType
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenExpiration)).Unix()

	return token
}

func (h *HTTP) getJWTConfig() (jwtCfg middleware.JWTConfig) {
	jwtCfg.SigningMethod = h.jwtSigningMethod.Name
	jwtCfg.SigningKey = []byte(h.jwtSigningKey)
	return
}
