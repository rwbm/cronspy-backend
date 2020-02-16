package transport

import (
	"cronspy/backend/pkg/api/job"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// DefaultPageSize configures the default number of records to return
	DefaultPageSize = 10
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

	jobs := e.Group("/jobs")
	jobs.GET("", h.userJobsHandler, IsUserLoggedIn) // get user jobs

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
	var page, pageSize int
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	if pageStr == "" {
		pageStr = "1"
	}
	if pageSizeStr == "" {
		pageSize = DefaultPageSize
	}

	page, errConv := strconv.Atoi(pageStr)
	if errConv != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidPage, ""))
	}

	if pageSize == 0 {
		pageSize, errConv = strconv.Atoi(pageSizeStr)
		if errConv != nil {
			return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidPageSize, ""))
		}
	}

	// get jobs
	jobs, p, err := h.svc.GetJobs(int(idUser), pageSize, page)
	if err != nil {
		return err
	}

	type response struct {
		Jobs       []model.Job      `json:"jobs,omitempty"`
		Pagination model.Pagination `json:"pagination,omitempty"`
	}

	return c.JSON(http.StatusOK, response{Jobs: jobs, Pagination: p})
}
