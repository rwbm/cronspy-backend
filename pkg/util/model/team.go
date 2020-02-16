package model

import "time"

// Team is a way of grouping users (only paid accounts)
type Team struct {
	ID          int       `gorm:"column:id_team;primary_key;AUTO_INCREMENT"`
	DateCreated time.Time `gorm:"NOT NULL"`
	Name        string    `gorm:"type:varchar(128);NOT NULL"`
}
