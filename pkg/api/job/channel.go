package job

import "cronspy/backend/pkg/util/model"

// SaveChannel saves a channel in the database
func (j *Job) SaveChannel(c *model.Channel) (err error) {

	err = j.database.SaveChannel(c)
	if err != nil {
		j.logger.Error("error saving channel", err, nil)
		return
	}

	return
}

// GetUserChannels returns the list of configured channels for a user
func (j *Job) GetUserChannels(idUser int) (channels []model.Channel, err error) {
	return
}
