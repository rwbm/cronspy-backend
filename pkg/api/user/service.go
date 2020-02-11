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
	RegisterUser(ec echo.Context, user *model.User) (err error)
}

type DB interface {
	RegisterUser(user *model.User) (id int, err error)
	GetUserByEmail(email string) (user model.User, err error)
}

// User defines the module for user related operations
type User struct {
	database DB
	logger   *log.Log
}

// creates new reseller service
func new(database DB, l *log.Log) *User {
	return &User{
		database: database,
		logger:   l,
	}
}

// Initialize initializes tax application service
func Initialize(ds *gorm.DB, l *log.Log) *User {
	return new(db.NewUserDB(ds), l)
}
