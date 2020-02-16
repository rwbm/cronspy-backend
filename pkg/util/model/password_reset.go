package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// PasswordReset represents a record used when the user forgets the password
type PasswordReset struct {
	ID            string    `gorm:"column:id_password_reset;primary_key"`
	DateCreated   time.Time `gorm:"NOT NULL"`
	DateUpdated   time.Time `gorm:"NOT NULL"`
	IDUser        int       `gorm:"NOT NULL"`
	LinkSentCount int       `gorm:"NOT NULL"`
	Validated     bool      `gorm:"DEFAULT(0);NOT NULL"`
	Used          bool      `gorm:"DEFAULT(0);NOT NULL"`
}

// TableName returns the table name for the model
func (PasswordReset) TableName() string {
	return "cronspy.password_resets"
}

// BeforeCreate sets the unique ID before record is saved in the database
func (pr *PasswordReset) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}
