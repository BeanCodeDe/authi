package db

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/BeanCodeDe/SpaceLight-Auth/internal/authErr"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	log "github.com/sirupsen/logrus"
)

type UserDB struct {
	ID        uuid.UUID `db:"id"`
	Password  string    `db:"password"`
	Roles     []string  `db:"roles"`
	CreatedOn time.Time `db:"created_on"`
	LastLogin time.Time `db:"last_login"`
}

func (user *UserDB) Create() error {
	log.Debugf("Create user")
	hash := getHash()

	if _, err := getConnection().Exec(context.Background(), "INSERT INTO spacelight.user(id, password,salt,roles,created_on,last_login) VALUES($1,MD5($2),$3,$4,$5,$6)", user.ID, user.Password+hash, hash, user.Roles, user.CreatedOn, user.LastLogin); err != nil {
		log.Errorf("Unknown error when inserting user: %v", err)
		return authErr.UnknownError
	}
	log.Debugf("User inserted into database")
	return nil
}

func getHash() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GetUserById(userId uuid.UUID) (*UserDB, error) {
	log.Debugf("Get user %s by UserId", userId)

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), getConnection(), &users, `SELECT id,roles::text[],created_on,last_login FROM spacelight.user WHERE id = $1`, userId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				log.Warnf("User with id %s not found", userId)
				return nil, authErr.UserNotFoundError
			}
		}
		log.Errorf("Unknown error when getting user by name: %v", err)
		return nil, authErr.UnknownError
	}

	if len(users) == 0 {
		log.Errorf("User with id %s not found.", userId)
		return nil, authErr.UserNotFoundError
	}

	if len(users) != 1 {
		log.Errorf("Cant find only one user. Len: %s, Userlist: %v", len(users), users)
		return nil, authErr.UnknownError
	}

	log.Debugf("Got user %v", users[0])
	return users[0], nil
}

func (user *UserDB) LoginUser() error {
	log.Debugf("Check password for user %s", user.ID)

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), getConnection(), &users, `SELECT id,roles::text[],created_on,last_login FROM spacelight.user WHERE id = $1 AND password = MD5(CONCAT($2::text,(SELECT salt FROM spacelight.user WHERE id = $1)::text))`, user.ID, user.Password); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				log.Warnf("No user with id %s and matching password found", user.ID)
				return authErr.WrongAuthDataError
			}
		}
		log.Errorf("Unknown error when checking user: %v", err)
		return authErr.UnknownError
	}

	if len(users) == 0 {
		log.Errorf("User with id %s not found. Probably wrong password", user.ID)
		return authErr.WrongAuthDataError
	}

	if len(users) != 1 {
		log.Errorf("Cant find only one user. Len: %s, Userlist: %v", len(users), users)
		return authErr.UnknownError
	}

	log.Debugf("Got user %v", users[0].ID)

	user.CreatedOn = users[0].CreatedOn
	user.LastLogin = users[0].LastLogin
	user.Roles = users[0].Roles

	return nil
}
