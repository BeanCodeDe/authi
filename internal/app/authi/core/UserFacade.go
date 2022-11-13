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
)

func (user *UserCore) Create() (err error) {

	creationTime := time.Now()
	user.CreatedOn = creationTime
	user.LastLogin = creationTime

	if err = user.mapToUserDB().Create(); err != nil {
		return err
	}

	return nil
}

func (user *UserCore) Login() (*TokenCore, error) {
	dbUser := user.mapToUserDB()
	err := dbUser.LoginUser()
	if err != nil {
		return nil, fmt.Errorf("something went wrong when logging in user, %v", user.ID)
	}
	return createJWTToken(user.ID)
}

func (user *UserCore) mapToUserDB() *db.UserDB {
	return &db.UserDB{ID: user.ID, Password: user.Password, CreatedOn: user.CreatedOn, LastLogin: user.LastLogin}
}
