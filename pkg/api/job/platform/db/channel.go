package db

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"fmt"

	"github.com/jinzhu/gorm"
)

// SaveChannel saves the channel on the database.
func (j *JobDB) SaveChannel(channel *model.Channel) (err error) {

	trx := j.ds.Begin()
	err = trx.Save(channel).Error

	if err == nil {

		// save config based on channel type
		switch channel.Type {

		case model.ChannelTypeEmail:
			cfg := channel.GetChannelEmail()
			cfg.ID = channel.ID
			err = trx.Save(&cfg).Error

		case model.ChannelTypeSlack:
			cfg := channel.GetChannelSlack()
			cfg.ID = channel.ID
			err = trx.Save(&cfg).Error

		case model.ChannelTypeWebHook:
			cfg := channel.GetChannelWebHook()
			cfg.ID = channel.ID
			err = trx.Save(&cfg).Error

		}

		if err != nil {
			trx.Rollback()
			return
		}

		// commit changes if everything was OK
		trx.Commit()
	} else {
		trx.Rollback()
	}

	return
}

// GetChannels returns the list of channels defined by a user;
// specific channel configuration can also be loaded.
func (j *JobDB) GetChannels(idUser int, loadChannelConfig bool) (channels []model.Channel, err error) {
	if err = j.ds.Model(model.Channel{}).Where("id_user = ?", idUser).Find(&channels).Error; err == nil {

		if loadChannelConfig {
			// load configuration for each channel
			for i := range channels {
				if err = j.getChannelConfig(&channels[i]); err != nil {
					break
				}
			}
		}
	}
	return
}

// GetChannel returns a channel by ID
func (j *JobDB) GetChannel(idChannel int, loadChannelConfig bool) (c model.Channel, err error) {

	err = j.ds.Model(c).Where("id_channel = ?", idChannel).First(&c).Error
	if err == nil {
		if loadChannelConfig {
			err = j.getChannelConfig(&c)
		}
	} else {
		if err == gorm.ErrRecordNotFound {
			err = exception.ErrRecordNotFound
		}
	}

	return
}

// DeleteChannel removes an existing channel and it's configuration from the database
func (j *JobDB) DeleteChannel(channel *model.Channel) (err error) {

	trx := j.ds.Begin()

	// delete channel config
	switch channel.Type {

	case model.ChannelTypeEmail:
		cfg := channel.GetChannelEmail()
		cfg.ID = channel.ID
		err = trx.Delete(&cfg).Error

	case model.ChannelTypeSlack:
		cfg := channel.GetChannelSlack()
		cfg.ID = channel.ID
		err = trx.Delete(&cfg).Error

	case model.ChannelTypeWebHook:
		cfg := channel.GetChannelWebHook()
		cfg.ID = channel.ID
		err = trx.Delete(&cfg).Error

	}

	if err != nil {
		trx.Rollback()
		return
	}

	// delete channel
	if err = trx.Delete(channel).Error; err != nil {
		trx.Rollback()
		return
	}

	// commit changes if everything was OK
	trx.Commit()
	return
}

// UpdateChannel saves configuration for an existing channel;
// only name and configuration can be changed.
func (j *JobDB) UpdateChannel(channel *model.Channel) (err error) {

	if channel.ID > 0 {
		trx := j.ds.Begin()

		if err = trx.Model(&channel).Updates(map[string]interface{}{"name": channel.Name}).Error; err == nil {
			// save configuration
			var cfg interface{}

			switch channel.Type {

			case model.ChannelTypeEmail:
				c := channel.GetChannelEmail()
				cfg = &c

			case model.ChannelTypeSlack:
				c := channel.GetChannelSlack()
				cfg = &c

			case model.ChannelTypeWebHook:
				c := channel.GetChannelWebHook()
				cfg = &c
			}

			if err = trx.Save(cfg).Error; err != nil {
				trx.Rollback()
				return
			}

			// save changes
			trx.Commit()

		} else {
			trx.Rollback()
		}
	}

	return
}

// saves channel configuration into the database
func (j *JobDB) saveCahnnelConfig(channel *model.Channel) (err error) {
	err = j.ds.Save(channel.Configuration).Error
	return
}

func (j *JobDB) getChannelConfig(c *model.Channel) (err error) {

	switch c.Type {

	case model.ChannelTypeEmail:
		m := &model.ChannelEmail{}
		q := j.ds.Model(m)
		if err = q.Where("id_channel = ?", c.ID).First(m).Error; err == nil {
			c.SetChannelEmail(*m)
		}

	case model.ChannelTypeWebHook:
		m := &model.ChannelWebHook{}
		q := j.ds.Model(m)
		if err = q.Where("id_channel = ?", c.ID).First(m).Error; err == nil {
			c.SetChannelWebHook(*m)
		}

	case model.ChannelTypeSlack:
		m := &model.ChannelSlack{}
		q := j.ds.Model(m)
		if err = q.Where("id_channel = ?", c.ID).First(m).Error; err == nil {
			c.SetChannelSlack(*m)
		}

	default:
		err = fmt.Errorf("channel type '%s' not supported", c.Type)
	}

	return
}
