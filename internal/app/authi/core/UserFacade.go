package core

import (
	"time"

	"github.com/BeanCodeDe/SpaceLight-Auth/internal/auth"
	"github.com/BeanCodeDe/SpaceLight-Auth/internal/db"
	"github.com/BeanCodeDe/SpaceLight-AuthMiddleware/authAdapter"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type (
	UserCore struct {
		ID        uuid.UUID
		Password  string
		Roles     []string
		CreatedOn time.Time
		LastLogin time.Time
	}
)

func (user *UserCore) Create() (createdUser *UserCore, err error) {
	log.Debugf("Create user")

	user.ID = uuid.New()

	creationTime := time.Now()
	user.CreatedOn = creationTime
	user.LastLogin = creationTime

	user.Roles = []string{authAdapter.UserRole}

	if err = user.mapToUserDB().Create(); err != nil {
		return nil, err
	}

	dbUser, err := db.GetUserById(user.ID)
	if err != nil {
		return nil, err
	}

	return mapToUserCore(dbUser), nil
}

func (user *UserCore) Login() (string, error) {
	log.Debugf("Login user %s", user.ID)
	dbUser := user.mapToUserDB()
	err := dbUser.LoginUser()
	if err != nil {
		log.Debugf("Something went wrong when logging in user, %v", user.ID)
		return "", err
	}
	coreUser := mapToUserCore(dbUser)
	return auth.CreateJWTToken(coreUser.ID, coreUser.Roles)
}

func (user *UserCore) mapToUserDB() *db.UserDB {
	return &db.UserDB{ID: user.ID, Password: user.Password, Roles: user.Roles, CreatedOn: user.CreatedOn, LastLogin: user.LastLogin}
}

func mapToUserCore(dbUser *db.UserDB) *UserCore {
	return &UserCore{ID: dbUser.ID, Password: dbUser.Password, Roles: dbUser.Roles, CreatedOn: dbUser.CreatedOn, LastLogin: dbUser.LastLogin}
}
