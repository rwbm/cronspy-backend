package db

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"

	"github.com/jinzhu/gorm"
)

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

// RegisterUser creates a new user in the database
func (c *UserDB) RegisterUser(user *model.User) (id int, err error) {
	err = c.ds.Create(&user).Error
	return
}

// GetUserByEmail finds a user by the email address
func (c *UserDB) GetUserByEmail(email string) (user model.User, err error) {
	if err = c.ds.Where("email = ?", email).First(&user).Error; err == gorm.ErrRecordNotFound {
		err = exception.ErrRecordNotFound
	}
	return
}
