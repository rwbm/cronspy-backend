package transport

import (
	"bytes"
	"cronspy/backend/pkg/api/user"
	"cronspy/backend/pkg/util/exception"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/model"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type loginResp struct {
	User        model.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

// ****************************************************
//
// DATABASE MOCK
//
// ****************************************************

func getDBMock(mockData bool) (db *DBMock) {
	db = &DBMock{}

	if mockData {
		// load some users
		u1 := &model.User{Email: "test.user.a1@cronspy.com", Name: "Test User A1", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u1.HashPassword()
		db.RegisterUser(u1)

		u2 := &model.User{Email: "test.user.a2@cronspy.com", Name: "Test User A2", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u2.HashPassword()
		db.RegisterUser(u2)

		u3 := &model.User{Email: "test.user.a3@cronspy.com", Name: "Test User A3", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u3.HashPassword()
		db.RegisterUser(u3)
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

//
// ============== USER REGISTRATION ==============

func runUserRegistration(user model.User) (resp model.User, err error) {
	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, true)

	// define request
	uJSON, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(uJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	if err = handler.userRegisterHandler(c); err == nil {
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
	}

	return
}

func TestUserRegistrationOK(t *testing.T) {
	// define user to test
	u := model.User{
		Email:    "test.user.1@cronspay.com",
		Name:     "Test User 1",
		Password: "abcd1234",
	}

	newUser, err := runUserRegistration(u)

	// assertions
	if assert.NoError(t, err) {
		// validate returned data
		assert.NotEqual(t, 0, newUser.ID)
		assert.Equal(t, "test.user.1@cronspay.com", newUser.Email)
		assert.Equal(t, "Test User 1", newUser.Name)
	}
}

func TestUserRegistrationInvalidEmail(t *testing.T) {
	// define user to test
	u := model.User{
		Email:    "test.user.1@invalidaddress",
		Name:     "Test User 1",
		Password: "abcd1234",
	}

	_, err := runUserRegistration(u)

	// assertions
	assert.Error(t, err)
}

func TestUserRegistrationInvalidPassword(t *testing.T) {
	// define user to test
	u := model.User{
		Email:    "test.user.1@cronspay.com",
		Name:     "Test User 1",
		Password: "bad",
	}

	_, err := runUserRegistration(u)

	// assertions
	assert.Error(t, err)
}

func TestUserRegistrationExistingUser(t *testing.T) {
	// define user to test
	u := model.User{
		Email:    "test.user.a1@cronspy.com",
		Name:     "Test User 1",
		Password: "anotherpassword",
	}

	_, err := runUserRegistration(u)

	// assertions
	assert.Error(t, err)
}

//
// ============== USER LOGIN ==============

func runUserLogin(username, password string) (resp loginResp, err error) {
	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, true)

	// define request
	payload := fmt.Sprintf(`{ "username":"%s", "password":"%s" }`, username, password)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	if err = handler.userLoginHandler(c); err == nil {
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
	}

	return
}

func TestUserLoginOK(t *testing.T) {
	r, err := runUserLogin("test.user.a1@cronspy.com", "abcd1234")

	// assertions
	if assert.NoError(t, err) {
		assert.Equal(t, "Test User A1", r.User.Name)
	}
}

func TestUserLoginWrongPassword(t *testing.T) {
	_, err := runUserLogin("test.user.a1@cronspy.com", "wrong-password")

	// assertions
	assert.Error(t, err)
}

func TestUserLoginInvalidEmail(t *testing.T) {
	_, err := runUserLogin("test.user.a1@invalid-host", "wrong-password")

	// assertions
	assert.Error(t, err)
}

func TestUserLoginInvalidPassword(t *testing.T) {
	_, err := runUserLogin("test.user.a1@cronspy.com", "bad")

	// assertions
	assert.Error(t, err)
}
