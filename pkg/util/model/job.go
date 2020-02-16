package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Job is a job configured for a user, to be monitored by the system
type Job struct {
	ID                      string     `gorm:"column:id_job;primary_key" json:"id"`
	IDUser                  int        `gorm:"NOT NULL" json:"id_user"`
	DateCreated             time.Time  `gorm:"NOT NULL" json:"date_created"`
	DateUpdated             time.Time  `gorm:"NOT NULL" json:"date_updated"`
	Name                    string     `gorm:"NOT NULL" json:"name"`
	JobType                 string     `gorm:"NOT NULL" json:"job_type"`
	Active                  bool       `gorm:"NOT NULL" json:"active"`
	Status                  string     `gorm:"NOT NULL" json:"status"`
	CronExpression          *string    `json:"cron_expression"`
	CronExpressionTimezone  *string    `json:"cron_expression_timezone"`
	DetectedIntervalMinutes *int       `json:"detected_interval"`
	Alerts                  []JobAlert `gorm:"foreignkey:ID" json:"-"`
}

func (Job) TableName() string {
	return "cronspy.jobs"
}

// BeforeCreate sets the unique ID before record is saved in the database
func (j *Job) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New().String())
	return nil
}

// JobAlert represent an alert definition, when something goes wrong
type JobAlert struct {
	ID                        int     `gorm:"column:id_alert;primary_key" json:"id"`
	IDJob                     int     `gorm:"NOT NULL" json:"id_job"`
	Target                    string  `gorm:"NOT NULL" json:"target"`
	MinutesBeforeNotification int     `gorm:"NOT NULL" json:"minutes_before_notification"`
	IDChannel                 int     `gorm:"NOT NULL" json:"id_channel"`
	Channel                   Channel `gorm:"foreignkey:IDChannel" json:"-"`
}

func (JobAlert) TableName() string {
	return "cronspy.job_alerts"
}
