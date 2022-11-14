package db

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type UserDB struct {
	ID        uuid.UUID `db:"id"`
	Password  string    `db:"password"`
	CreatedOn time.Time `db:"created_on"`
	LastLogin time.Time `db:"last_login"`
}

func (user *UserDB) Create() error {
	hash := getHash()
	if _, err := getConnection().Exec(context.Background(), "INSERT INTO auth.user(id, password,salt,created_on,last_login) VALUES($1,MD5($2),$3,$4,$5)", user.ID, user.Password+hash, hash, user.CreatedOn, user.LastLogin); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return errormessages.UserAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting user: %v", err)
	}
	return nil
}

func UpdateRefreshToken(userId uuid.UUID, refreshToken string) error {
	if _, err := getConnection().Exec(context.Background(), "UPDATE auth.user SET refresh_token=$1 WHERE id=$2", refreshToken, userId); err != nil {
		return fmt.Errorf("unknown error when updating refresh token of user %s user: %v", userId, err)
	}
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
	var users []*UserDB
	if err := pgxscan.Select(context.Background(), getConnection(), &users, `SELECT id,created_on,last_login FROM auth.user WHERE id = $1`, userId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				return nil, fmt.Errorf("user with id %s not found", userId)
			}
		}
		return nil, fmt.Errorf("unknown error when getting user by name: %v", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with id %s not found", userId)
	}

	if len(users) != 1 {
		return nil, fmt.Errorf("cant find only one user. Len: %v, Userlist: %v", len(users), users)
	}

	return users[0], nil
}

func (user *UserDB) LoginUser() error {

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), getConnection(), &users, `SELECT id,created_on,last_login FROM auth.user WHERE id = $1 AND password = MD5(CONCAT($2::text,(SELECT salt FROM auth.user WHERE id = $1)::text))`, user.ID, user.Password); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				return fmt.Errorf("no user with id %s and matching password found", user.ID)
			}
		}
		return fmt.Errorf("unknown error when checking user: %v", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("user with id %s not found. Probably wrong password", user.ID)
	}

	if len(users) != 1 {
		return fmt.Errorf("cant find only one user. Len: %v, Userlist: %v", len(users), users)
	}

	user.CreatedOn = users[0].CreatedOn
	user.LastLogin = users[0].LastLogin

	return nil
}

func CheckRefreshToken(userId uuid.UUID, refreshToken string) error {

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), getConnection(), &users, `SELECT id FROM auth.user WHERE id = $1 AND refresh_token = $2`, userId, refreshToken); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				return fmt.Errorf("no user with id %s and matching refresh token found", userId)
			}
		}
		return fmt.Errorf("unknown error when checking user: %v", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("user with id %s not found. Probably wrong refresh token", userId)
	}

	if len(users) != 1 {
		return fmt.Errorf("cant find only one user. Len: %v, Userlist: %v", len(users), users)
	}

	return nil
}
