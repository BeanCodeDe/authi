package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
)

var (
	//go:embed migration/sqlite/*.up.sql
	sqLiteMigrationFs embed.FS
)

type (
	SqlLiteConnection struct {
		db *sql.DB
	}
)

func newSqlLiteConnectionConnection() (Connection, error) {
	dbFilePath := util.GetEnvWithFallback("SQL_LITE_PATH", "./data.db")
	info, err := os.Stat(dbFilePath)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error while checking state of db: %w", err)
	}

	if err != nil && errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(dbFilePath)
		if err != nil {
			return nil, fmt.Errorf("error while creating db file in %s: %w", dbFilePath, err)
		}
	}

	if info.IsDir() {
		return nil, fmt.Errorf("db file [%s] is not a file but a directory", dbFilePath)
	}

	err = migrateSqLiteDatabase(dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("error while creating db file in %s: %w", dbFilePath, err)
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open sql file: %v", err)
	}
	return &SqlLiteConnection{db}, nil
}

func migrateSqLiteDatabase(url string) error {
	d, err := iofs.New(sqLiteMigrationFs, "migration/postgres")
	if err != nil {
		return fmt.Errorf("error while creating instance of migration scrips: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, url)
	if err != nil {
		return fmt.Errorf("error while creating instance of migration scrips: %w", err)
	}
	err = m.Up()
	if err != nil {
		return fmt.Errorf("error while migrating: %w", err)
	}
	return nil
}

func (connection *SqlLiteConnection) Close() {
	connection.db.Close()
}

func (connection *SqlLiteConnection) CreateUser(user *UserDB, hash string) error {
	stmt, err := connection.db.Prepare("INSERT INTO auth.user(id, password,salt,created_on,last_login) VALUES($1,MD5($2),$3,$4,$5)")
	if err != nil {
		return fmt.Errorf("unknown error when creating insert statement: %v", err)
	}
	_, err = stmt.Exec(user.ID, user.Password+hash, hash, user.CreatedOn, user.LastLogin)
	if err != nil {
		return fmt.Errorf("unknown error when inserting user: %v", err)
	}
	return nil
}

func (connection *SqlLiteConnection) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	stmt, err := connection.db.Prepare("UPDATE auth.user SET refresh_token=$1, refresh_token_expire=$2 WHERE id=$3")
	if err != nil {
		return fmt.Errorf("unknown error when creating update refresh token statement: %v", err)
	}
	_, err = stmt.Exec(refreshToken, refreshTokenExpireAt, userId)
	if err != nil {
		return fmt.Errorf("unknown error when updating refresh token: %v", err)
	}
	return nil
}

func (connection *SqlLiteConnection) LoginUser(user *UserDB) error {
	rows, err := connection.db.Query("SELECT id,created_on,last_login FROM auth.user WHERE id = $1 AND password = MD5(CONCAT($2::text,(SELECT salt FROM auth.user WHERE id = $1)::text))", user.ID, user.Password)
	if err != nil {
		return fmt.Errorf("unknown error when checking user: %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("user with id %s not found. Probably wrong password", user.ID)
	}

	if rows.Next() {
		return fmt.Errorf("cant find only one user")
	}

	return nil
}

func (connection *SqlLiteConnection) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {
	rows, err := connection.db.Query("SELECT id FROM auth.user WHERE id = $1 AND refresh_token = $2 AND refresh_token_expire > now()", userId, refreshToken)
	if err != nil {
		return fmt.Errorf("unknown error when checking user: %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("user with id %s not found. Probably wrong refresh token", userId)
	}

	if rows.Next() {
		return fmt.Errorf("cant find only one user")
	}

	return nil
}

func (connection *SqlLiteConnection) UpdatePassword(userId uuid.UUID, password string, hash string) error {
	stmt, err := connection.db.Prepare("UPDATE auth.user SET password=MD5($1), salt=$2 WHERE id=$3")
	if err != nil {
		return fmt.Errorf("unknown error when creating update password statement: %v", err)
	}
	_, err = stmt.Exec(password+hash, hash, userId)
	if err != nil {
		return fmt.Errorf("unknown error when updating password of user %s error: %v", userId, err)
	}
	return nil
}

func (connection *SqlLiteConnection) DeleteUser(userId uuid.UUID) error {
	stmt, err := connection.db.Prepare("DELETE FROM auth.user WHERE id=$1")
	if err != nil {
		return fmt.Errorf("unknown error when creating delete statement: %v", err)
	}
	_, err = stmt.Exec(userId)
	if err != nil {
		return fmt.Errorf("unknown error when deleting user %s error: %v", userId, err)
	}
	return nil
}
