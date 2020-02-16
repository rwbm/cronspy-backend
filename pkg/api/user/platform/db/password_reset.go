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
func (c *UserDB) GetPasswordResetByID(id string, trx *gorm.DB) (reset model.PasswordReset, err error) {

	ds := c.ds
	if trx != nil {
		ds = trx

		defer func() {
			if r := recover(); r != nil {
				ds.Rollback()
			}
		}()
	}

	if err = ds.Where("id_password_reset = ?", id).First(&reset).Error; err == gorm.ErrRecordNotFound {
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

// UpdatePasswordResetCount updates the sent count
func (c *UserDB) UpdatePasswordResetCount(id string, countValue int) (err error) {
	return c.ds.Model(model.PasswordReset{}).Where("id_password_reset = ?", id).Update("link_sent_count", countValue).Error
}

// ValidatePasswordReset updates the validated field with true
func (c *UserDB) ValidatePasswordReset(id string) (err error) {
	return c.ds.Model(model.PasswordReset{}).Where("id_password_reset = ?", id).Update("validated", true).Error
}

// MarkPasswordResetAsUsed updates the used field with true
func (c *UserDB) MarkPasswordResetAsUsed(id string, trx *gorm.DB) (err error) {
	ds := c.ds
	if trx != nil {
		ds = trx

		defer func() {
			if r := recover(); r != nil {
				ds.Rollback()
			}
		}()
	}

	return ds.Model(model.PasswordReset{}).Where("id_password_reset = ?", id).Update("used", true).Error
}
