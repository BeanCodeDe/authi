package core

import (
	"errors"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type (
	ConnectionMock struct {
		closeRecordArray              []*closeRecord
		createUserRecordArray         []*createUserRecord
		updateRefreshTokenRecordArray []*updateRefreshTokenRecord
		loginUserRecordArray          []*loginUserRecord
		checkRefreshTokenRecordArray  []*checkRefreshTokenRecord
		updatePasswordRecordArray     []*updatePasswordRecord
		deleteUserRecordArray         []*deleteUserRecord
		createUserReturn              error
		updateRefreshTokenReturn      error
		loginUserReturn               error
		checkRefreshReturn            error
		updatePasswordReturn          error
		deleteUserReturn              error
	}

	closeRecord struct {
	}

	createUserRecord struct {
		user *db.UserDB
		hash string
	}

	updateRefreshTokenRecord struct {
		userId               uuid.UUID
		refreshToken         string
		refreshTokenExpireAt time.Time
	}

	loginUserRecord struct {
		user *db.UserDB
	}
	checkRefreshTokenRecord struct {
		userId       uuid.UUID
		refreshToken string
	}
	updatePasswordRecord struct {
		userId   uuid.UUID
		password string
		hash     string
	}

	deleteUserRecord struct {
		userId uuid.UUID
	}
)

var (
	userId         = uuid.New()
	refreshToken   = "someRefreshToken"
	authenticate   = &adapter.AuthenticateDTO{}
	privateKeyPath = "../../../../deployments/token/privat/jwtRS256.key"
	errUnknown     = errors.New("some error")
)

func TestRandomString(t *testing.T) {
	randomString := randomString()
	assert.Len(t, randomString, 32)
	assert.Regexp(t, "[a-zA-Z0-9]*", randomString)
}

func (connection *ConnectionMock) Close() {
	closeRecord := &closeRecord{}
	connection.closeRecordArray = append(connection.closeRecordArray, closeRecord)
}

func (connection *ConnectionMock) CreateUser(user *db.UserDB, hash string) error {
	createUserRecord := &createUserRecord{user: user, hash: hash}
	connection.createUserRecordArray = append(connection.createUserRecordArray, createUserRecord)
	return connection.createUserReturn
}

func (connection *ConnectionMock) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	updateRefreshToken := &updateRefreshTokenRecord{userId: userId, refreshToken: refreshToken, refreshTokenExpireAt: refreshTokenExpireAt}
	connection.updateRefreshTokenRecordArray = append(connection.updateRefreshTokenRecordArray, updateRefreshToken)
	return connection.updateRefreshTokenReturn
}

func (connection *ConnectionMock) LoginUser(user *db.UserDB) error {
	loginUserRecord := &loginUserRecord{user: user}
	connection.loginUserRecordArray = append(connection.loginUserRecordArray, loginUserRecord)
	return connection.loginUserReturn
}

func (connection *ConnectionMock) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {
	checkRefreshTokenRecord := &checkRefreshTokenRecord{userId: userId, refreshToken: refreshToken}
	connection.checkRefreshTokenRecordArray = append(connection.checkRefreshTokenRecordArray, checkRefreshTokenRecord)
	return connection.checkRefreshReturn
}

func (connection *ConnectionMock) UpdatePassword(userId uuid.UUID, password string, hash string) error {
	updatePasswordRecord := &updatePasswordRecord{userId: userId, password: password, hash: hash}
	connection.updatePasswordRecordArray = append(connection.updatePasswordRecordArray, updatePasswordRecord)
	return connection.updatePasswordReturn
}

func (connection *ConnectionMock) DeleteUser(userId uuid.UUID) error {
	deleteUserRecord := &deleteUserRecord{userId: userId}
	connection.deleteUserRecordArray = append(connection.deleteUserRecordArray, deleteUserRecord)
	return connection.deleteUserReturn
}
