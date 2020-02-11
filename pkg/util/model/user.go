package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account types
const (
	AccountTypeFree     = "FREE"
	AccountTypeStartup  = "STARTUP"
	AccountTypeBusiness = "BUSINESS"
	AccountTypeCustom   = "CUSTOM"
)

// User represents a customer user
type User struct {
	ID             int       `gorm:"column:id_user;primary_key;AUTO_INCREMENT" json:"id"`
	DateCreated    time.Time `gorm:"NOT NULL" json:"date_created"`
	DateUpdated    time.Time `gorm:"NOT NULL" json:"date_updated"`
	Email          string    `gorm:"type:varchar(128);unique_index;NOT NULL" json:"email"`
	Name           string    `gorm:"type:varchar(128);NOT NULL" json:"name"`
	HashedPassword string    `gorm:"column:password;type:varchar(128);NOT NULL" json:"-"`
	Password       string    `gorm:"-" json:"password,omitempty"`
	AccountType    string    `gorm:"type:default(FREE);NOT NULL" json:"account_type"`
}

// HashPassword takes the value from .Password and stores
// the BCrypted version of if in HashedPassword
func (u *User) HashPassword() {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	u.HashedPassword = string(hashedBytes)
}

// ValidatePassword validated `password` against the stored hashed password
func (u *User) ValidatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

// CleanPassword removes password values from the model
func (u *User) CleanPassword() {
	u.Password = ""
	u.HashedPassword = ""
}

// Team is a way of grouping users (only paid accounts)
type Team struct {
	ID          int       `gorm:"column:id_team;primary_key;AUTO_INCREMENT"`
	DateCreated time.Time `gorm:"NOT NULL"`
	Name        string    `gorm:"type:varchar(128);NOT NULL"`
}
