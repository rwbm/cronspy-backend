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
	IDUser        int                    `gorm:"NOT NULL" json:"-"`
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
		cwh.ID = c.ID
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

// SetChannelWebHook sets the configuration map with
// the values of the model
func (c *Channel) SetChannelWebHook(cwh ChannelWebHook) {
	if c.Configuration == nil {
		c.Configuration = make(map[string]interface{})
	}

	if cwh.BaseURL == "" {
		c.Configuration["base_url"] = cwh.BaseURL
	}
	if cwh.PayloadType == "" {
		c.Configuration["payload_type"] = cwh.PayloadType
	}
	if cwh.BasicAuthUsername != nil {
		c.Configuration["basic_auth_username"] = *cwh.BasicAuthUsername
	}
	if cwh.BasicAuthPassword != nil {
		c.Configuration["basic_auth_password"] = *cwh.BasicAuthPassword
	}
}

// GetChannelSlack returns Configuration as `ChannelSlack`
func (c *Channel) GetChannelSlack() (cs ChannelSlack) {
	if c.Configuration != nil {
		cs.ID = c.ID
		if v, ok := c.Configuration["base_url"].(string); ok {
			cs.BaseURL = v
		}
		if v, ok := c.Configuration["slack_channel_name"].(string); ok {
			cs.SlackChannelName = &v
		}
	}
	return
}

// SetChannelSlack sets the configuration map with
// the values of the model
func (c *Channel) SetChannelSlack(cwh ChannelSlack) {
	if c.Configuration == nil {
		c.Configuration = make(map[string]interface{})
	}

	if cwh.BaseURL == "" {
		c.Configuration["base_url"] = cwh.BaseURL
	}
	if cwh.SlackChannelName != nil {
		c.Configuration["slack_channel_name"] = *cwh.SlackChannelName
	}
}

// GetChannelEmail returns Configuration as `ChannelEmail`
func (c *Channel) GetChannelEmail() (ce ChannelEmail) {
	if c.Configuration != nil {
		ce.ID = c.ID
		if v, ok := c.Configuration["email"].(string); ok {
			ce.Email = v
		}
	}
	return
}

// SetChannelEmail sets the configuration map with
// the values of the model
func (c *Channel) SetChannelEmail(cwh ChannelEmail) {
	if c.Configuration == nil {
		c.Configuration = make(map[string]interface{})
	}

	if cwh.Email != "" {
		c.Configuration["email"] = cwh.Email
	}
}

// ChannelEmail contains the configuration when a channel
// is of type `ChannelEmail`
type ChannelEmail struct {
	ID    int    `gorm:"column:id_channel;primary_key" json:"-"`
	Email string `gorm:"NOT NULL" json:"email"`
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
