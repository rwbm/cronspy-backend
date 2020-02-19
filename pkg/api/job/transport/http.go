package transport

import (
	"cronspy/backend/pkg/api/job"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// DefaultPageSize configures the default number of records to return
	DefaultPageSize = 15
	// DefaultJobName contains a default name for jobs that are created without one
	DefaultJobName = "Job Monitor"
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

	// configure routes
	jobs := e.Group("/jobs")
	jobs.GET("", h.userJobsHandler, IsUserLoggedIn)       // get user jobs
	jobs.POST("", h.createJobHandler, IsUserLoggedIn)     // create job
	jobs.GET("/:job-id", h.getJobHandler, IsUserLoggedIn) // get job by id

	channels := e.Group("/channels")
	channels.POST("", h.createChannelHandler, IsUserLoggedIn)               // create channel
	channels.DELETE("/:channel-id", h.deleteChannelHandler, IsUserLoggedIn) // delete channel
	channels.PUT("/:channel-id", h.updateChannelHandler, IsUserLoggedIn)    // update channel

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
	idUser, _, err := h.getUserID(c)
	if err != nil {
		return err
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

//
// --- GET JOB ---
//
func (h *HTTP) getJobHandler(c echo.Context) error {
	// get user id
	idUser, _, err := h.getUserID(c)
	if err != nil {
		return err
	}

	// get jobs
	job, err := h.svc.GetJob(c.Param("job-id"))
	if err != nil {
		return err
	}

	if job.IDUser != int(idUser) {
		return echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
	}

	return c.JSON(http.StatusOK, job)
}

//
// --- CREATE JOB ---
//
func (h *HTTP) createJobHandler(c echo.Context) error {
	// get user id
	idUser, _, err := h.getUserID(c)
	if err != nil {
		return err
	}

	payload := new(model.Job)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// validate input
	if fields := h.validateCreateJobInput(payload); fields != "" {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMapWithFields(exception.CodeInvalidFields, "", fields))
	}

	payload.IDUser = idUser
	if err := h.svc.SaveJob(payload); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, payload)
}

//
// --- CREATE CHANNEL ---
//
func (h *HTTP) createChannelHandler(c echo.Context) error {
	// get user id
	idUser, emailAddress, err := h.getUserID(c)
	if err != nil {
		return err
	}

	payload := new(model.Channel)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// validate input
	if fields := h.validateCreateChannelInput(payload); fields != "" {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMapWithFields(exception.CodeInvalidFields, "", fields))
	}

	// if no email address was provided, use the user's
	if payload.Type == model.ChannelTypeEmail {
		addEmail := true
		if payload.Configuration != nil {
			if _, found := payload.Configuration["email"]; found {
				addEmail = false
			}
		} else {
			payload.Configuration = make(map[string]interface{})
		}

		if addEmail {
			payload.Configuration["email"] = emailAddress
		}
	}

	// save channel
	payload.IDUser = idUser
	if err := h.svc.SaveChannel(payload); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, payload)
}

//
// --- DELETE CHANNEL ---
//
func (h *HTTP) deleteChannelHandler(c echo.Context) error {

	// get user id
	idUser, _, err := h.getUserID(c)
	if err != nil {
		return err
	}

	// get channel id from path
	idChannelStr := c.Param("channel-id")
	idChannel, errConv := strconv.Atoi(idChannelStr)
	if errConv != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidEntityID, errConv.Error()))
	}

	// save channel
	errDelete := h.svc.DeleteChannel(idChannel, idUser)
	if errDelete != nil {
		return errDelete
	}

	return c.NoContent(http.StatusOK)
}

//
// --- UPDATE CHANNEL ---
//
func (h *HTTP) updateChannelHandler(c echo.Context) error {

	// get user id
	idUser, _, err := h.getUserID(c)
	if err != nil {
		return err
	}

	// get channel id from path
	idChannelStr := c.Param("channel-id")
	idChannel, errConv := strconv.Atoi(idChannelStr)
	if errConv != nil {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeInvalidEntityID, errConv.Error()))
	}

	payload := new(model.Channel)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// save channel
	errUpdate := h.svc.UpdateChannel(idChannel, idUser, payload)
	if errUpdate != nil {
		return errUpdate
	}

	return c.NoContent(http.StatusOK)
}

//
// --- private methods ---
//

// get user ID and email from request context (must be authenticated)
func (h *HTTP) getUserID(c echo.Context) (id int, email string, err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	idUser, okID := claims["id"].(float64)
	emailStr, okEmail := claims["email"].(string)

	if !okID || !okEmail {
		err = echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
		return
	}

	id = int(idUser)
	email = emailStr

	return
}

// validate create job fields
func (h *HTTP) validateCreateJobInput(j *model.Job) (fields string) {
	invalidFields := []string{}

	// for cons, we need a con expression
	if j.JobType == model.JobTypeCron {
		if j.CronExpression == nil {
			// TODO: validate cron expression too
			invalidFields = append(invalidFields, "cron_expression")
		}
		if j.CronExpressionTimezone == nil {
			// TODO: validate timezone
			invalidFields = append(invalidFields, "cron_expression_timezone")
		}
	}

	if j.Name == "" {
		j.Name = DefaultJobName
	}

	if len(invalidFields) > 0 {
		fields = strings.Join(invalidFields, ",")
	}

	return
}

// validate create channel fields
func (h *HTTP) validateCreateChannelInput(c *model.Channel) (fields string) {
	invalidFields := []string{}

	if c.Type != model.ChannelTypeEmail && c.Type != model.ChannelTypeSlack && c.Type != model.ChannelTypeWebHook {
		invalidFields = append(invalidFields, "type")
	}

	// configuration is mandatory for some types
	switch c.Type {
	case model.ChannelTypeSlack, model.ChannelTypeWebHook:
		if c.Configuration == nil {
			invalidFields = append(invalidFields, "configuration")
		}
	}

	if len(invalidFields) > 0 {
		fields = strings.Join(invalidFields, ",")
	}

	return
}
