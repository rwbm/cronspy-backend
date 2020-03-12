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
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// helper types

type loginResp struct {
	User        model.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

var (
	passwordRestDummyID1, passwordRestDummyID2 string
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
		u1 := &model.User{Email: "test.user.a1@cronspy.com", Name: "Test User A1", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u1.HashPassword()
		db.RegisterUser(u1)

		u2 := &model.User{Email: "test.user.a2@cronspy.com", Name: "Test User A2", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u2.HashPassword()
		db.RegisterUser(u2)

		u3 := &model.User{Email: "test.user.a3@cronspy.com", Name: "Test User A3", Password: "abcd1234", DateCreated: time.Now(), DateUpdated: time.Now(), AccountType: model.AccountTypeFree}
		u3.HashPassword()
		db.RegisterUser(u3)

		// password resets
		pr1 := &model.PasswordReset{ID: "dummy-uuid-1", IDUser: 2, LinkSentCount: 1, DateCreated: time.Now(), DateUpdated: time.Now()}
		db.CreatePasswordReset(pr1)
		passwordRestDummyID1 = pr1.ID

		pr2 := &model.PasswordReset{ID: "dummy-uuid-2", IDUser: 3, LinkSentCount: 1, DateCreated: time.Now().Add(-1 * time.Hour), DateUpdated: time.Now().Add(-1 * time.Hour)}
		db.CreatePasswordReset(pr2)
		passwordRestDummyID2 = pr2.ID
	}

	return
}

type DBMock struct {
	users          []model.User
	passwordResets []model.PasswordReset

	currentUserID int
	mux           sync.Mutex
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
	reset.ID = uuid.New().String()
	db.passwordResets = append(db.passwordResets, *reset)
	return
}

func (db *DBMock) GetPasswordResetByID(id string, trx *gorm.DB) (reset model.PasswordReset, err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].ID == id {
				reset = db.passwordResets[i]
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) GetPasswordResetByUser(idUser int) (reset model.PasswordReset, err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].IDUser == idUser {
				reset = db.passwordResets[i]
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) DeletePasswordReset(id string) (err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].ID == id {
				db.passwordResets = append(db.passwordResets[:i], db.passwordResets[i+1:]...)
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) UpdatePasswordResetCount(id string, countValue int) (err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].ID == id {
				db.passwordResets[i].LinkSentCount = countValue
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) ValidatePasswordReset(id string) (err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].ID == id {
				db.passwordResets[i].Validated = true
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
	return
}

func (db *DBMock) MarkPasswordResetAsUsed(id string, trx *gorm.DB) (err error) {
	if len(db.passwordResets) > 0 {
		for i := range db.passwordResets {
			if db.passwordResets[i].ID == id {
				db.passwordResets[i].Used = true
				return
			}
		}
	}

	err = exception.ErrRecordNotFound
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

//
// ============== PASSWORD RESET ==============

func runPasswordReset(email string) (id string, err error) {
	// create server and handler
	e := echo.New()
	handler := getHTTPHandler(e, true)

	// define request
	payload := fmt.Sprintf(`{ "email":"%s" }`, email)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// call handler
	if err = handler.userPasswordResetRequestHandler(c); err == nil {
		resp := make(map[string]interface{})
		if err = json.Unmarshal(rec.Body.Bytes(), &resp); err == nil {
			id = resp["id"].(string)
		}

	}

	return
}

func TestPasswordResetOK(t *testing.T) {
	id, err := runPasswordReset("test.user.a1@cronspy.com")

	// assertions
	if assert.NoError(t, err) {
		assert.NotEmpty(t, id)
	}
}

func TestPasswordResetInvalidUser(t *testing.T) {
	_, err := runPasswordReset("unknown-user@cronspy.com")

	// assertions
	assert.Error(t, err)
}

func TestPasswordResetExistingError(t *testing.T) {
	_, err := runPasswordReset("test.user.a2@cronspy.com")

	// assertions
	assert.Error(t, err)
}

func TestPasswordResetExistingOK(t *testing.T) {
	id, err := runPasswordReset("test.user.a3@cronspy.com")

	// assertions
	if assert.NoError(t, err) {
		assert.Equal(t, passwordRestDummyID2, id)
	}
}
