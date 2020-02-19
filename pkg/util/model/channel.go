package model

// Channel type definitions
const (
	ChannelTypeEmail   = "EMAIL"
	ChannelTypeWebHook = "WEB_HOOK"
	ChannelTypeSlack   = "SLACK"
)

// Channel represents a notification channel
type Channel struct {
	ID            int                    `gorm:"column:id_channel;primary_key" json:"id"`
	IDUser        int                    `gorm:"NOT NULL" json:"id_user"`
	Type          string                 `gorm:"NOT NULL" json:"type"`
	Name          string                 `gorm:"NOT NULL" json:"name"`
	Configuration map[string]interface{} `gorm:"-" json:"configuration"`
}

// TableName returns the table name for the model
func (Channel) TableName() string {
	return "cronspy.channels"
}

// GetChannelWebHook returns Configuration as `ChannelWebHook`
func (c *Channel) GetChannelWebHook() (cwh ChannelWebHook) {
	if c.Configuration != nil {
		if v, ok := c.Configuration["base_url"].(string); ok {
			cwh.BaseURL = v
		}
		if v, ok := c.Configuration["payload_type"].(string); ok {
			cwh.PayloadType = v
		}
		if v, ok := c.Configuration["basic_auth_username"].(string); ok {
			cwh.BasicAuthUsername = &v
		}
		if v, ok := c.Configuration["basic_auth_password"].(string); ok {
			cwh.BasicAuthPassword = &v
		}
	}
	return
}

// GetChannelSlack returns Configuration as `ChannelSlack`
func (c *Channel) GetChannelSlack() (cs ChannelSlack) {
	if c.Configuration != nil {
		if v, ok := c.Configuration["base_url"].(string); ok {
			cs.BaseURL = v
		}
		if v, ok := c.Configuration["slack_channel_name"].(string); ok {
			cs.SlackChannelName = &v
		}
	}
	return
}

// GetChannelEmail returns Configuration as `ChannelEmail`
func (c *Channel) GetChannelEmail() (ce ChannelEmail) {
	if c.Configuration != nil {
		if v, ok := c.Configuration["email_address"].(string); ok {
			ce.Address = v
		}
	}
	return
}

// ChannelEmail contains the configuration when a channel
// is of type `ChannelEmail`
type ChannelEmail struct {
	ID      int    `gorm:"column:id_channel;primary_key" json:"-"`
	Address string `gorm:"NOT NULL" json:"email_address"`
}

// TableName returns the table name for the model
func (ChannelEmail) TableName() string {
	return "cronspy.channels_email"
}

// ChannelWebHook contains the configuration when a channel
// is of type `ChannelTypeWebHook`
type ChannelWebHook struct {
	ID                int     `gorm:"column:id_channel;primary_key" json:"id"`
	BaseURL           string  `gorm:"NOT NULL" json:"base_url"`
	PayloadType       string  `gorm:"NOT NULL" json:"payload_type"`
	BasicAuthUsername *string `json:"basic_auth_username"`
	BasicAuthPassword *string `json:"basic_auth_password"`
}

// TableName returns the table name for the model
func (ChannelWebHook) TableName() string {
	return "cronspy.channels_webhook"
}

// ChannelSlack contains the configuration when a channel
// is of type `ChannelTypeSlack`
type ChannelSlack struct {
	ID               int     `gorm:"column:id_channel;primary_key" json:"id"`
	BaseURL          string  `gorm:"NOT NULL" json:"base_url"`
	SlackChannelName *string `json:"slack_channel_name"`
}

// TableName returns the table name for the model
func (ChannelSlack) TableName() string {
	return "cronspy.channels_slack"
}
