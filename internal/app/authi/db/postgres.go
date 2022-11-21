package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	PostgresConnection struct {
		dbPool *pgxpool.Pool
	}
)

func NewPostgresConnection() (*PostgresConnection, error) {

	user := config.PostgresUser
	name := config.PostgresDB
	password := config.PostgresPassword
	host := config.PostgresHost
	port := config.PostgresPort

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)
	var err error
	dbPool, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return &PostgresConnection{dbPool: dbPool}, nil
}

func (connection *PostgresConnection) Close() {
	connection.dbPool.Close()
}

func (connection *PostgresConnection) CreateUser(user *UserDB, hash string) error {
	if _, err := connection.dbPool.Exec(context.Background(), "INSERT INTO auth.user(id, password,salt,created_on,last_login) VALUES($1,MD5($2),$3,$4,$5)", user.ID, user.Password+hash, hash, user.CreatedOn, user.LastLogin); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return errormessages.ErrUserAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting user: %v", err)
	}
	return nil
}

func (connection *PostgresConnection) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	if _, err := connection.dbPool.Exec(context.Background(), "UPDATE auth.user SET refresh_token=$1, refresh_token_expire=$2 WHERE id=$3", refreshToken, refreshTokenExpireAt, userId); err != nil {
		return fmt.Errorf("unknown error when updating refresh token of user %s user: %v", userId, err)
	}
	return nil
}

func (connection *PostgresConnection) GetUserById(userId uuid.UUID) (*UserDB, error) {
	var users []*UserDB
	if err := pgxscan.Select(context.Background(), connection.dbPool, &users, `SELECT id,created_on,last_login FROM auth.user WHERE id = $1`, userId); err != nil {
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

func (connection *PostgresConnection) LoginUser(user *UserDB) error {

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), connection.dbPool, &users, `SELECT id,created_on,last_login FROM auth.user WHERE id = $1 AND password = MD5(CONCAT($2::text,(SELECT salt FROM auth.user WHERE id = $1)::text))`, user.ID, user.Password); err != nil {
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

func (connection *PostgresConnection) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {

	var users []*UserDB
	if err := pgxscan.Select(context.Background(), connection.dbPool, &users, `SELECT id FROM auth.user WHERE id = $1 AND refresh_token = $2 AND refresh_token_expire > now()`, userId, refreshToken); err != nil {
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
