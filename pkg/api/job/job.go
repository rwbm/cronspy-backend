package job

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetJobs returns the list of configured monitors for a user
func (j *Job) GetJobs(idUser int, pageSize, page int) (jobs []model.Job, p model.Pagination, err error) {
	jobs, p, err = j.database.GetJobs(idUser, pageSize, page)
	if err != nil {
		j.logger.Error("error loading user jobs", err, map[string]interface{}{"id_user": idUser})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		return
	}

	return
}

// GetJob return a job data by the ID
func (j *Job) GetJob(id string) (job model.Job, err error) {
	job, err = j.database.GetJobByID(id)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusNotFound, exception.GetErrorMap(exception.CodeNotFound, ""))
		} else {
			j.logger.Error("error loading job", err, map[string]interface{}{"id_job": id})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		}
	}
	return
}

// SaveJob saves a new job in the database
func (j *Job) SaveJob(job *model.Job) (err error) {

	err = j.database.SaveJob(job)
	if err != nil {
		j.logger.Error("error saving job", err, map[string]interface{}{"id_user": job.IDUser})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
	}

	return
}
