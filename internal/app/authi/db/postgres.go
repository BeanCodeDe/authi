package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

var (
	//go:embed migration/postgres/*.up.sql
	postgresMigrationFs embed.FS
)

type (
	postgresConnection struct {
		dbPool *pgxpool.Pool
	}
)

func newPostgresConnection() (Connection, error) {
	user := util.GetEnvWithFallback("POSTGRES_USER", "postgres")
	dbName := util.GetEnvWithFallback("POSTGRES_DB", "postgres")
	password, err := util.GetEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("postgres password has to be set: %w", err)
	}
	host := util.GetEnvWithFallback("POSTGRES_HOST", "postgres")
	port, err := util.GetEnvIntWithFallback("POSTGRES_PORT", 5432)
	options := util.GetEnvWithFallback("POSTGRES_OPTIONS", "sslmode=disable")
	migrationOptions := util.GetEnvWithFallback("POSTGRES_MIGRATION_OPTIONS", "&x-migrations-table=authi-migration")

	if err != nil {
		return nil, fmt.Errorf("port is not a number: %w", err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", user, password, host, port, dbName, options)
	err = migratePostgresDatabase(url + migrationOptions)
	if err != nil {
		return nil, fmt.Errorf("error while migrating database: %w", err)
	}

	dbPool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return &postgresConnection{dbPool: dbPool}, nil
}

func (connection *postgresConnection) Close() {
	connection.dbPool.Close()
}

func migratePostgresDatabase(url string) error {
	d, err := iofs.New(postgresMigrationFs, "migration/postgres")
	if err != nil {
		return fmt.Errorf("error while creating instance of migration scrips: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, url)
	if err != nil {
		return fmt.Errorf("error while creating instance of migration scrips: %w", err)
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("error while migrating: %w", err)
	}
	return nil
}

func (connection *postgresConnection) CreateUser(user *UserDB, hash string) error {
	if _, err := connection.dbPool.Exec(context.Background(), "INSERT INTO auth.user(id, password,salt,created_on,last_login, init_user) VALUES($1,MD5($2),$3,$4,$5,$6)", user.ID, user.Password+hash, hash, user.CreatedOn, user.LastLogin, user.InitUser); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrUserAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting user: %v", err)
	}
	return nil
}

func (connection *postgresConnection) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	if _, err := connection.dbPool.Exec(context.Background(), "UPDATE auth.user SET refresh_token=$1, refresh_token_expire=$2 WHERE id=$3", refreshToken, refreshTokenExpireAt, userId); err != nil {
		return fmt.Errorf("unknown error when updating refresh token of user %s error: %v", userId, err)
	}
	return nil
}

func (connection *postgresConnection) LoginUser(user *UserDB) error {

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

func (connection *postgresConnection) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {

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

func (connection *postgresConnection) UpdatePassword(userId uuid.UUID, password string, hash string) error {
	if _, err := connection.dbPool.Exec(context.Background(), "UPDATE auth.user SET password=MD5($1), salt=$2 WHERE id=$3", password+hash, hash, userId); err != nil {
		return fmt.Errorf("unknown error when updating password of user %s error: %v", userId, err)
	}
	return nil
}

func (connection *postgresConnection) DeleteUser(userId uuid.UUID) error {
	if _, err := connection.dbPool.Exec(context.Background(), "DELETE FROM auth.user WHERE id=$1", userId); err != nil {
		return fmt.Errorf("unknown error when deleting user %s error: %v", userId, err)
	}
	return nil
}

func (connection *postgresConnection) DeleteInitUsers() error {
	if _, err := connection.dbPool.Exec(context.Background(), "DELETE FROM auth.user WHERE init_user=TRUE"); err != nil {
		return fmt.Errorf("unknown error when deleting init users: %v", err)
	}
	return nil
}
