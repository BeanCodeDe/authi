package db

import (
	"errors"
	"strings"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/util"
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
		UpdatePassword(userId uuid.UUID, password string, hash string) error
		DeleteUser(userId uuid.UUID) error
	}
)

func NewConnection() (Connection, error) {
	switch db := strings.ToLower(util.GetEnvWithFallback("DATABASE", "SQLite")); db {
	case "sqlite":
		return newSqlLiteConnectionConnection()
	case "postgresql":
		return newPostgresConnection()
	default:
		return nil, errors.New("no configuration for %s found")
	}
}
