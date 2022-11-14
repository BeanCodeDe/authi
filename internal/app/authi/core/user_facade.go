package core

import (
	"fmt"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/google/uuid"
)

type (
	UserCore struct {
		ID        uuid.UUID
		Password  string
		CreatedOn time.Time
		LastLogin time.Time
	}

	User interface {
		Create() error
		Login() (*TokenCore, error)
		GetId() uuid.UUID
		SetId(uuid.UUID)
		GetCreatedOn() time.Time
		GetLastLogin() time.Time
	}
)

func (user *UserCore) Create() error {

	creationTime := time.Now()
	user.CreatedOn = creationTime
	user.LastLogin = creationTime

	if err := user.mapToUserDB().Create(); err != nil {
		return err
	}

	return nil
}

func (user *UserCore) Login() (*TokenCore, error) {
	dbUser := user.mapToUserDB()
	err := dbUser.LoginUser()
	if err != nil {
		return nil, fmt.Errorf("something went wrong when logging in user, %v: %v", user.ID, err)
	}
	return createJWTToken(user.ID)
}

func (user *UserCore) GetId() uuid.UUID {
	return user.ID
}

func (user *UserCore) SetId(userId uuid.UUID) {
	user.ID = userId
}

func (user *UserCore) GetCreatedOn() time.Time {
	return user.CreatedOn
}

func (user *UserCore) GetLastLogin() time.Time {
	return user.LastLogin
}

func (user *UserCore) mapToUserDB() *db.UserDB {
	return &db.UserDB{ID: user.ID, Password: user.Password, CreatedOn: user.CreatedOn, LastLogin: user.LastLogin}
}
