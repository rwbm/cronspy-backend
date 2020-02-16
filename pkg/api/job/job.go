package job

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"

	"github.com/labstack/echo"
)

// GetJobs returns the list of configured monitors for a user
func (j *Job) GetJobs(idUser int, count, offset int) (jobs []model.Job, err error) {

	if count == 0 {
		count = DefaultPageSize
	}

	jobs, err = j.database.GetJobs(idUser, count, offset)
	if err != nil {
		j.logger.Error("error loading user jobs", err, map[string]interface{}{"id_user": idUser})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		return
	}

	return
}
