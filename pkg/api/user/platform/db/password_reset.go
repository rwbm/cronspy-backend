package db

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"

	"github.com/jinzhu/gorm"
)

// CreatePasswordReset creates a password reset record
func (c *UserDB) CreatePasswordReset(reset *model.PasswordReset) (err error) {
	err = c.ds.Create(&reset).Error
	return
}

// GetPasswordResetByID finds an existing password reset by ID
func (c *UserDB) GetPasswordResetByID(id string) (reset model.PasswordReset, err error) {
	if err = c.ds.Where("id_password_reset = ?", id).First(&reset).Error; err == gorm.ErrRecordNotFound {
		err = exception.ErrRecordNotFound
	}
	return
}

// GetPasswordResetByUser finds an existing password reset by the user ID
func (c *UserDB) GetPasswordResetByUser(idUser int) (reset model.PasswordReset, err error) {
	if err = c.ds.Where("id_user = ?", idUser).First(&reset).Error; err == gorm.ErrRecordNotFound {
		err = exception.ErrRecordNotFound
	}
	return
}

// DeletePasswordReset deletes an existing password reset record
func (c *UserDB) DeletePasswordReset(id string) error {
	return c.ds.Where("id_password_reset = ?", id).Delete(model.PasswordReset{}).Error
}

// UpdatePasswordResetCount updated the sent count
func (c *UserDB) UpdatePasswordResetCount(id string, countValue int) (err error) {
	return c.ds.Model(model.PasswordReset{}).Where("id_password_reset = ?", id).Update("link_sent_count", countValue).Error
}
