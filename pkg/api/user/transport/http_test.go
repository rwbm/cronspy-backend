package transport

import (
	"bytes"
	"cronspy/backend/pkg/api/user"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/model"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ****************************************************
//
// DATABASE MOCK
//
// ****************************************************

func getDBMock(mockData bool) (db *DBMock) {
	db = &DBMock{}

	if mockData {
		// load some users
		db.RegisterUser(&model.User{Email: "test.user.a1@cronspy.com", Name: "Test User A1", Password: "asdasd"})
		db.RegisterUser(&model.User{Email: "test.user.a2@cronspy.com", Name: "Test User A2", Password: "asdasd"})
		db.RegisterUser(&model.User{Email: "test.user.a3@cronspy.com", Name: "Test User A3", Password: "asdasd"})
	}

	return
}

type DBMock struct {
	users         []model.User
	passwordReset []model.PasswordReset

	currentUserID          int
	currentPasswordResetID int
	mux                    sync.Mutex
}

func (db *DBMock) Transaction() *gorm.DB {
	return &gorm.DB{}
}

func (db *DBMock) RegisterUser(user *model.User) (id int, err error) {
	db.mux.Lock()
	db.currentUserID++
	user.ID = db.currentUserID
	db.mux.Unlock()

	db.users = append(db.users, *user)
	return
}

func (db *DBMock) GetUserByEmail(email string) (user model.User, err error) {
	if len(db.users) > 0 {
		for i := range db.users {
			if db.users[i].Email == email {
				user = db.users[i]
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) GetUserByID(idUser int) (user model.User, err error) {
	if len(db.users) > 0 {
		for i := range db.users {
			if db.users[i].ID == idUser {
				user = db.users[i]
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) UpdateUserPassword(idUser int, newPassword string, trx *gorm.DB) (err error) {
	return
}

func (db *DBMock) UpdateUser(user *model.User, fields map[string]interface{}) (err error) {
	return
}

func (db *DBMock) CreatePasswordReset(reset *model.PasswordReset) (err error) {
	return
}

func (db *DBMock) GetPasswordResetByID(id string, trx *gorm.DB) (reset model.PasswordReset, err error) {
	return
}

func (db *DBMock) GetPasswordResetByUser(idUser int) (reset model.PasswordReset, err error) {
	return
}

func (db *DBMock) DeletePasswordReset(id string) (err error) {
	return
}

func (db *DBMock) UpdatePasswordResetCount(id string, countValue int) (err error) {
	return
}

func (db *DBMock) ValidatePasswordReset(id string) (err error) {
	return
}

func (db *DBMock) MarkPasswordResetAsUsed(id string, trx *gorm.DB) (err error) {
	return
}

// get HTTP handler
func getHTTPHandler(e *echo.Echo, mockData bool) (h HTTP) {

	mockDB := getDBMock(mockData)
	logger := log.New()
	userService := user.Initialize(nil, mockDB, logger, 5)

	h = NewHTTP(userService, "myTestingKey", jwt.SigningMethodHS512, e)
	return
}

func TestUserRegistrationOK(t *testing.T) {

	// define user to test
	u := model.User{
		Email:    "test.user.1@cronspay.com",
		Name:     "Test User 1",
		Password: "abcd1234",
	}
	uJSON, _ := json.Marshal(u)

	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, true)

	// define request
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(uJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	err := handler.userRegisterHandler(c)

	// assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		// TODO: capturar JSON y validar datos retornados

		// assert.Equal(t, userJSON, rec.Body.String())
	}
}

func TestUserRegistrationInvalidEmail(t *testing.T) {

	// define user to test
	u := model.User{
		Email:    "test.user.1@invalidaddress",
		Name:     "Test User 1",
		Password: "abcd1234",
	}
	uJSON, _ := json.Marshal(u)

	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, false)

	// define request
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(uJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	err := handler.userRegisterHandler(c)

	// assertions
	assert.Error(t, err)
}

func TestUserRegistrationInvalidPassword(t *testing.T) {

	// define user to test
	u := model.User{
		Email:    "test.user.1@cronspay.com",
		Name:     "Test User 1",
		Password: "abc123",
	}
	uJSON, _ := json.Marshal(u)

	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, false)

	// define request
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(uJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	err := handler.userRegisterHandler(c)

	// assertions
	assert.Error(t, err)
}

func TestUserRegistrationExsitingUser(t *testing.T) {

	// define user to test
	u := model.User{
		Email:    "test.user.a1@cronspy.com",
		Name:     "Test User 1",
		Password: "abc123",
	}
	uJSON, _ := json.Marshal(u)

	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, false)

	// define request
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(uJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	err := handler.userRegisterHandler(c)

	// assertions
	assert.Error(t, err)
}
