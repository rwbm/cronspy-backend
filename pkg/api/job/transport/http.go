package transport

import (
	"cronspy/backend/pkg/api/job"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"

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
	svc              job.Service
	jwtSigningKey    string
	jwtSigningMethod *jwt.SigningMethodHMAC
}

// NewHTTP creates new http service to handle request to /user
func NewHTTP(svc job.Service, jwtSigningKey string, jwtSigningMethod *jwt.SigningMethodHMAC, e *echo.Echo) {
	h := HTTP{
		svc:              svc,
		jwtSigningKey:    jwtSigningKey,
		jwtSigningMethod: jwtSigningMethod,
	}

	// define logged user check function
	IsUserLoggedIn = middleware.JWTWithConfig(h.getJWTConfig())

	job := e.Group("/jobs")

	// --- Auth required ---
	job.GET("", h.userJobsHandler, IsUserLoggedIn)
}

func (h *HTTP) getJWTConfig() (jwtCfg middleware.JWTConfig) {
	jwtCfg.SigningMethod = h.jwtSigningMethod.Name
	jwtCfg.SigningKey = []byte(h.jwtSigningKey)
	return
}

//
// --- GET USER JOBS ---.
//
func (h *HTTP) userJobsHandler(c echo.Context) error {

	// get user id
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idUser, ok := claims["id"].(float64)

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
	}

	// get pagination data
	// pageStr := c.QueryParam("page")
	// pageCountStr := c.QueryParam("page_count")

	// get jobs
	jobs, err := h.svc.GetJobs(int(idUser), 0, 0)
	if err != nil {
		return err
	}

	type response struct {
		Jobs []model.Job `json:"jobs"`
	}

	return c.JSON(http.StatusOK, response{Jobs: jobs})
}
