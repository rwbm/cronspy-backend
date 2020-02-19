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

// SaveCahnnelConfig saves channel configuration into the database
func (j *JobDB) saveCahnnelConfig(channel *model.Channel) (err error) {
	err = j.ds.Save(channel.Configuration).Error
	return
}

// GetChannels returns the list of channels defined by a user;
// specific channel configuration can also be loaded.
func (j *JobDB) GetChannels(idUser int, loadChannelConfig bool) (channels []model.Channel, err error) {
	if err = j.ds.Model(model.Channel{}).Where("id_user = ?", idUser).Find(&channels).Error; err == nil {

		if loadChannelConfig {

			for i := range channels {
				var q *gorm.DB

				switch channels[i].Type {

				case model.ChannelTypeEmail:
					q = j.ds.Model(model.ChannelEmail{})

				case model.ChannelTypeWebHook:
					q = j.ds.Model(model.ChannelWebHook{})

				case model.ChannelTypeSlack:
					q = j.ds.Model(model.ChannelSlack{})

				default:
					err = fmt.Errorf("channel type '%s' not supported", channels[i].Type)
					return
				}

				if err = q.Where("id_channel = ?", channels[i].ID).First(&channels[i].Configuration).Error; err != nil {
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
