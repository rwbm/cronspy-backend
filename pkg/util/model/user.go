package model

import "time"

// User represents a customer user
type User struct {
	ID          int       `gorm:"primary_key;AUTO_INCREMENT"`
	DateCreated time.Time `gorm:"NOT NULL"`
	DateUpdated time.Time `gorm:"NOT NULL"`
	Email       string    `gorm:"type:varchar(128);unique_index;NOT NULL"`
	Name        string    `gorm:"type:varchar(128);NOT NULL"`
	Password    string    `gorm:"type:varchar(128);NOT NULL"`
}
