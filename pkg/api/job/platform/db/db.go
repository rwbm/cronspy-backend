package db

import "github.com/jinzhu/gorm"

// NewJobDB returns a new tax database instance
func NewJobDB(ds *gorm.DB) (c *JobDB) {
	c = new(JobDB)
	c.ds = ds
	return
}

// JobDB contains the services to handle jobs
type JobDB struct {
	ds *gorm.DB
}

// Transaction returns a new database transaction
func (c *JobDB) Transaction() *gorm.DB {
	return c.ds.Begin()
}
