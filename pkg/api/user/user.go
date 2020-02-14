package user

import (
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// GetJWTExpiration returns the configured token expiration
func (u *User) GetJWTExpiration() int {
	return u.tokenExpiration
}

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

// Login handles a user login request, by loading the user from the DB
// and validating the password hash
func (u *User) Login(username, password string) (user model.User, err error) {

	// get user
	user, err = u.database.GetUserByEmail(username)
	if err == nil {

		// check password
		if !user.ValidatePassword(password) {
			err = echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
		} else {
			user.CleanPassword()
		}

	} else {
		if err == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		}
	}

	return
}

// ChangePassword handles logic for password changes
func (u *User) ChangePassword(idUser int, oldPassword, newPassword string) (err error) {

	// get user
	user, err := u.database.GetUserByID(idUser)
	if err == nil {

		// check password
		if !user.ValidatePassword(oldPassword) {
			err = echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
		} else {
			// update password
			if errUpdate := u.database.UpdateUserPassword(user.ID, newPassword); errUpdate != nil {
				u.logger.Error("error updating user password", errUpdate, map[string]interface{}{"id_user": idUser})
				err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errUpdate.Error()))
			}
		}

	} else {
		if err == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusUnauthorized, exception.GetErrorMap(exception.CodeInvalidPassword, ""))
		} else {
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, err.Error()))
		}
	}

	return
}

// ResetPassword holds the logic for password reset.
func (u *User) ResetPassword(email string) (resetID string, err error) {

	// check if users exists
	user, errGetUser := u.database.GetUserByEmail(email)
	if errGetUser != nil {
		if errGetUser == exception.ErrRecordNotFound {
			u.logger.Warn("a password reset operation was sent for a NON existing user", map[string]interface{}{"email": email})
			err = echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeUnknownUser, ""))

			//
			// HACK: otra opcion es retornar OK y no dar informacion si la cuenta de email existe o no
			//

		} else {
			u.logger.Error("error loading user by email from database", errGetUser, map[string]interface{}{"email": email})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errGetUser.Error()))
		}

		return
	}

	// check if there's a password reset already created for this user
	reset, ok, errReset := u.getOrCreatePasswordReset(user.ID)
	if errReset != nil {
		u.logger.Error("error creating password reset for user", errReset, map[string]interface{}{"id_user": user.ID})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errReset.Error()))
		return
	}

	// wait MinutesToWaitBeforeEmailResend before re-sending the email
	if !ok {
		err = echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeNeedToWaitBeforeResend, ""))
		return
	}

	// check if max was reached
	if reset.LinkSentCount > MaxPasswordRests {
		err = echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodeMaxPasswordResetReached, ""))
		return
	}

	// TODO: send email

	resetID = reset.ID

	return
}

// ValidateResetPassword is invoked when the user clicks the reset password URL
// provided by email; if not found, just return 404
func (u *User) ValidateResetPassword(resetID string) (err error) {

	// find password reset token
	reset, errGetReset := u.database.GetPasswordResetByID(resetID)
	if errGetReset != nil {
		if errGetReset == exception.ErrRecordNotFound {
			err = echo.NewHTTPError(http.StatusNotFound, exception.GetErrorMap(exception.CodeNotFound, ""))
		} else {
			u.logger.Error("error loading password reset from database", errGetReset, map[string]interface{}{"id_reset": resetID})
			err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errGetReset.Error()))
		}
		return
	}

	// check if EXPIRED
	if time.Since(reset.DateUpdated).Hours() > 24 {
		return echo.NewHTTPError(http.StatusBadRequest, exception.GetErrorMap(exception.CodePasswordResetTokenExpired, ""))
	}

	// mark password reset as vaidated
	if errUpdate := u.database.ValidatePasswordReset(resetID); errUpdate != nil {
		u.logger.Error("error updating password reset", errUpdate, map[string]interface{}{"id_reset": resetID})
		err = echo.NewHTTPError(http.StatusInternalServerError, exception.GetErrorMap(exception.CodeInternalServerError, errUpdate.Error()))
	}

	return
}

// looks for an existing password reset for the user; if it doesn't exist a new record is crated
func (u *User) getOrCreatePasswordReset(idUser int) (r model.PasswordReset, ok bool, err error) {

	r, err = u.database.GetPasswordResetByUser(idUser)
	if err != nil {
		if err == exception.ErrRecordNotFound {
			// create a new one
			r.IDUser = idUser
			r.DateCreated = time.Now()
			r.DateUpdated = time.Now()
			r.LinkSentCount = 1
			r.Validated = false
			r.Used = false

			err = u.database.CreatePasswordReset(&r)
			ok = true
		}
	} else {
		// wait MinutesToWaitBeforeEmailResend
		if time.Since(r.DateUpdated).Minutes() > (time.Duration(MinutesToWaitBeforeEmailResend) * time.Minute).Minutes() {
			r.LinkSentCount++
			err = u.database.UpdatePasswordResetCount(r.ID, r.LinkSentCount)
			ok = true
		}
	}

	return
}
