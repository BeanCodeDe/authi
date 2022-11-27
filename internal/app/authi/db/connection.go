package db

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserDB struct {
		ID        uuid.UUID `db:"id"`
		Password  string    `db:"password"`
		CreatedOn time.Time `db:"created_on"`
		LastLogin time.Time `db:"last_login"`
	}
	Connection interface {
		Close()
		CreateUser(user *UserDB, hash string) error
		UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error
		LoginUser(user *UserDB) error
		CheckRefreshToken(userId uuid.UUID, refreshToken string) error
	}
)
