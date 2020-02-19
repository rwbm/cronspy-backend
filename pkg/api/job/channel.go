package job

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

// DeleteChannel handles channel deletion
func (j *Job) DeleteChannel(idChannel, idUser int) (err error) {

	// get channel
	c, err := j.database.GetChannel(idChannel, true)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusNotFound, exception.GetErrorMap(exception.CodeNotFound, ""))
		} else {
			j.logger.Error("error loading channel", err, map[string]interface{}{"id_channel": idChannel, "id_user": idUser})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		}
		return
	}

	// check if the user is the owner
	if c.IDUser != idUser {
		err = echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
		return
	}

	// delete channel
	if errDelete := j.database.DeleteChannel(&c); errDelete != nil {
		j.logger.Error("error deleting channel", errDelete, map[string]interface{}{"id_channel": idChannel})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errDelete.Error()))
	}

	return
}

// UpdateChannel handles channel updates
func (j *Job) UpdateChannel(idChannel, idUser int, channel *model.Channel) (err error) {

	// get channel
	c, err := j.database.GetChannel(idChannel, true)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusNotFound, exception.GetErrorMap(exception.CodeNotFound, ""))
		} else {
			j.logger.Error("error loading channel", err, map[string]interface{}{"id_channel": idChannel, "id_user": idUser})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		}
		return
	}

	// check if the user is the owner
	if c.IDUser != idUser {
		err = echo.NewHTTPError(http.StatusForbidden, exception.GetErrorMap(exception.CodeUnauthorized, ""))
		return
	}

	// update channel
	c.Name = channel.Name
	c.Configuration = channel.Configuration

	if errUpdate := j.database.UpdateChannel(&c); errUpdate != nil {
		j.logger.Error("error updating channel data", errUpdate, map[string]interface{}{"id_channel": idChannel})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errUpdate.Error()))
	}

	return
}
