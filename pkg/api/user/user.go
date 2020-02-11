package user

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RegisterUser holds the logic to create a new user in the database
func (u *User) RegisterUser(ec echo.Context, user *model.User) (err error) {

	// check if the user already exists
	_, errGet := u.database.GetUserByEmail(user.Email)
	if errGet == nil {

		// user exists
		err = echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeUserExists, ""))

	} else {
		if errGet != exception.ErrRecordNotFound {
			u.logger.Error("error loading user by email", errGet, map[string]interface{}{"email": user.Email})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errGet.Error()))
			return
		}

		// create user
		user.DateCreated = time.Now()
		user.DateUpdated = time.Now()
		user.HashPassword()
		user.AccountType = model.AccountTypeFree // register users are always FREE

		if _, errSave := u.database.RegisterUser(user); errSave != nil {
			u.logger.Error("error creating user", errSave, nil)
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errSave.Error()))
		}

		// clean passwords for security
		user.CleanPassword()
	}

	return
}
