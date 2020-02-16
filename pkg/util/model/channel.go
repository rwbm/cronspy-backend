package model

// Channel type definitions
const (
	ChannelTypeEmail   = "EMAIL"
	ChannelTypeWebHook = "WEB_HOOK"
	ChannelTypeSlack   = "SLACK"
)

// Channel represents a notification channel
type Channel struct {
	ID     int    `gorm:"column:id_channel;primary_key" json:"id"`
	IDUser int    `gorm:"NOT NULL" json:"id_user"`
	Type   string `gorm:"NOT NULL" json:"type"`
	Name   string `gorm:"NOT NULL" json:"name"`
}

func (Channel) TableName() string {
	return "cronspy.channels"
}
