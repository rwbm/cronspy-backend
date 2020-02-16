package user

import (
	"cronspy/backend/pkg/api/user/platform/db"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// Service holds the functions delcared in the service interface
type Service interface {
	GetJWTExpiration() int
	RegisterUser(ec echo.Context, user *model.User) (err error)
	Login(username, password string) (user model.User, err error)
	ChangePassword(idUser int, oldPassword, newPassword string) (err error)
	ResetPassword(email string) (resetID string, err error)
	ValidateResetPassword(resetID string) (err error)
	ChangePasswordWithReset(resetToken, newPassword string) (err error)
}

// DB holds the functions for database access
type DB interface {
	Transaction() *gorm.DB

	// User
	RegisterUser(user *model.User) (id int, err error)
	GetUserByEmail(email string) (user model.User, err error)
	GetUserByID(idUser int) (user model.User, err error)
	UpdateUserPassword(idUser int, newPassword string, trx *gorm.DB) (err error)

	// Password resets
	CreatePasswordReset(reset *model.PasswordReset) (err error)
	GetPasswordResetByID(id string, trx *gorm.DB) (reset model.PasswordReset, err error)
	GetPasswordResetByUser(idUser int) (reset model.PasswordReset, err error)
	DeletePasswordReset(id string) error
	UpdatePasswordResetCount(id string, countValue int) (err error)
	ValidatePasswordReset(id string) (err error)
	MarkPasswordResetAsUsed(id string, trx *gorm.DB) (err error)
}

// User defines the module for user related operations
type User struct {
	database        DB
	logger          *log.Log
	tokenExpiration int
}

// creates new reseller service
func new(database DB, l *log.Log, tokenExpiration int) *User {
	return &User{
		database:        database,
		logger:          l,
		tokenExpiration: tokenExpiration,
	}
}

// Initialize initializes tax application service
func Initialize(ds *gorm.DB, l *log.Log, tokenExpiration int) *User {
	return new(db.NewUserDB(ds), l, tokenExpiration)
}
