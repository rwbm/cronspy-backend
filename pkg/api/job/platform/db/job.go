package db

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// GetJobs returns the list of jobs for a user
func (j *JobDB) GetJobs(idUser int, pageSize, page int) (jobs []model.Job, p model.Pagination, err error) {

	offset := 0
	if page > 1 {
		offset = ((page - 1) * pageSize)
	}

	// get total records
	totalRecords := 0
	if err = j.ds.Model(model.Job{}).Where("id_user = ?", idUser).Count(&totalRecords).Error; err != nil {
		return
	}

	// if page requested is invalid, return no results
	if page > totalRecords/pageSize {
		return
	}

	// get jobs
	q := j.ds.Model(model.Job{}).Where("id_user = ?", idUser).Order("date_created asc")
	q = q.Offset(offset).Limit(pageSize)
	err = q.Find(&jobs).Error

	if err == nil {
		p.Page = page
		p.PageSize = pageSize
		p.TotalRows = totalRecords
	}

	return
}

// GetJobByID return a job data by the ID
func (j *JobDB) GetJobByID(id string) (job model.Job, err error) {
	if err = j.ds.Model(model.Job{}).Where("id_job = ?", id).First(&job).Error; err == gorm.ErrRecordNotFound {
		err = exception.ErrRecordNotFound
	}
	return
}

// SaveJob saves a job in the database
func (j *JobDB) SaveJob(job *model.Job) (err error) {

	if job != nil {
		job.DateCreated = time.Now()
		job.DateUpdated = time.Now()
		job.Active = true
		job.Status = model.JobStatusUnknown
	} else {
		return errors.New("invalid entity")
	}

	err = j.ds.Save(job).Error
	return
}
