package db

import "github.com/jinzhu/gorm"

// NewUserDB returns a new tax database instance
func NewUserDB(ds *gorm.DB) (c *UserDB) {
	c = new(UserDB)
	c.ds = ds
	return
}

// UserDB contains the services to handle users
type UserDB struct {
	ds *gorm.DB
}

// Transaction returns a new database transaction
func (c *UserDB) Transaction() *gorm.DB {
	return c.ds.Begin()
}
