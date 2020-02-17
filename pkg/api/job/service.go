package job

import (
	"cronspy/backend/pkg/api/job/platform/db"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/model"

	"github.com/jinzhu/gorm"
)

// Service holds the functions delcared in the service interface
type Service interface {
	GetJobs(idUser int, pageSize, page int) (jobs []model.Job, p model.Pagination, err error)
	GetJob(id string) (job model.Job, err error)
	CreateJob(job *model.Job) (err error)
}

// DB holds the functions for database access
type DB interface {
	Transaction() *gorm.DB

	GetJobs(idUser int, count, offset int) (jobs []model.Job, p model.Pagination, err error)
	GetJobByID(id string) (job model.Job, err error)
	CreateJob(job *model.Job) (err error)
}

// Job defines the module for user related operations
type Job struct {
	database DB
	logger   *log.Log
}

// creates new reseller service
func new(database DB, l *log.Log) *Job {
	return &Job{
		database: database,
		logger:   l,
	}
}

// Initialize initializes tax application service
func Initialize(ds *gorm.DB, l *log.Log) *Job {
	return new(db.NewJobDB(ds), l)
}
