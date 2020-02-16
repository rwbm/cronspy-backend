package db

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

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

// GetUserByID finds a user by ID
func (c *UserDB) GetUserByID(idUser int) (user model.User, err error) {
	if err = c.ds.Where("id_user = ?", idUser).First(&user).Error; err == gorm.ErrRecordNotFound {
		err = exception.ErrRecordNotFound
	}
	return
}

// UpdateUser updates the list of fields for a user
func (c *UserDB) UpdateUser(user *model.User, fields map[string]interface{}) (err error) {
	return c.ds.Model(user).Update(fields).Error
}

// UpdateUserPassword updated the user with the new password hash
func (c *UserDB) UpdateUserPassword(idUser int, newPassword string, trx *gorm.DB) (err error) {

	ds := c.ds
	if trx != nil {
		ds = trx

		defer func() {
			if r := recover(); r != nil {
				ds.Rollback()
			}
		}()
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10); err == nil {
		err = ds.Model(model.User{}).Update("password", []byte(hashedPassword)).Error
	}

	return
}
