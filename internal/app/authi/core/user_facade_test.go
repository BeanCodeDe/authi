package core

import (
	"errors"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type (
	ConnectionMock struct {
		t                             *testing.T
		closeRecordArray              []*closeRecord
		createUserRecordArray         []*createUserRecord
		updateRefreshTokenRecordArray []*updateRefreshTokenRecord
		loginUserRecordArray          []*loginUserRecord
		checkRefreshTokenRecordArray  []*checkRefreshTokenRecord
		createUserReturn              error
		updateRefreshTokenReturn      error
		loginUserReturn               error
		checkRefreshReturn            error
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
)

var (
	userId         = uuid.New()
	refreshToken   = "someRefreshToken"
	authenticate   = &AuthenticateDTO{}
	privateKeyPath = "../../../../deployments/token/privat/jwtRS256.key"
)

// CreateUser Test

func TestCreateUser_Successfully(t *testing.T) {
	dbConnection := &ConnectionMock{t: t}
	userFacade := &UserFacade{dbConnection: dbConnection}

	if assert.NoError(t, userFacade.CreateUser(userId, authenticate)) {
		assert.Equal(t, 0, len(dbConnection.closeRecordArray))
		assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
		assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
		assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
		assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))

		assert.NotNil(t, dbConnection.createUserRecordArray[0].user)
		assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
		assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
	}
}

func TestCreateUser_CreateUser_UnknownError(t *testing.T) {
	unknownError := errors.New("some error")
	dbConnection := &ConnectionMock{t: t, createUserReturn: unknownError}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, authenticate)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))

	assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
	assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
}

func TestCreateUser_CreateUser_AlreadyExists(t *testing.T) {
	dbConnection := &ConnectionMock{t: t, createUserReturn: errormessages.ErrUserAlreadyExists}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, authenticate)
	assert.Equal(t, errormessages.ErrUserAlreadyExists, err)
	assert.Equal(t, 1, len(dbConnection.createUserRecordArray))

	assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
	assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
}

// RefreshToken Test

func TestRefreshToken_Successfully(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	dbConnection := &ConnectionMock{t: t}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))

	assert.Equal(t, userId, dbConnection.checkRefreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, dbConnection.checkRefreshTokenRecordArray[0].refreshToken)

	assert.Equal(t, userId, dbConnection.updateRefreshTokenRecordArray[0].userId)
	assert.Len(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt)

	assert.Regexp(t, "([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]*)", tokenResponseDTO.AccessToken)
	assert.Less(t, int(time.Now().Unix()), tokenResponseDTO.ExpiresIn)
	assert.Equal(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, tokenResponseDTO.RefreshToken)
	assert.Equal(t, int(dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt.Unix()), tokenResponseDTO.RefreshExpiresIn)
}

func TestRefreshToken_CheckRefreshToken_UnknownError(t *testing.T) {
	unknownError := errors.New("some error")
	dbConnection := &ConnectionMock{t: t, checkRefreshReturn: unknownError}

	userFacade := &UserFacade{dbConnection: dbConnection}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))

	assert.Equal(t, userId, dbConnection.checkRefreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, dbConnection.checkRefreshTokenRecordArray[0].refreshToken)

	assert.Nil(t, tokenResponseDTO)
	assert.NotNil(t, err)
}

func TestRefreshToken_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	unknownError := errors.New("some error")
	dbConnection := &ConnectionMock{t: t, updateRefreshTokenReturn: unknownError}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))

	assert.Equal(t, userId, dbConnection.checkRefreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, dbConnection.checkRefreshTokenRecordArray[0].refreshToken)

	assert.Equal(t, userId, dbConnection.updateRefreshTokenRecordArray[0].userId)
	assert.Len(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt)

	assert.Nil(t, tokenResponseDTO)
	assert.NotNil(t, err)
}

// LoginUser Test

func TestLoginUser_Successfully(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	dbConnection := &ConnectionMock{t: t}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, authenticate)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))

	assert.NotNil(t, dbConnection.loginUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.loginUserRecordArray[0].user.ID)

	assert.Equal(t, userId, dbConnection.updateRefreshTokenRecordArray[0].userId)
	assert.Len(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt)

	assert.Regexp(t, "([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]*)", tokenResponseDTO.AccessToken)
	assert.Less(t, int(time.Now().Unix()), tokenResponseDTO.ExpiresIn)
	assert.Equal(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, tokenResponseDTO.RefreshToken)
	assert.Equal(t, int(dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt.Unix()), tokenResponseDTO.RefreshExpiresIn)

}

func TestLoginUser_LoginUser_UnknownError(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	unknownError := errors.New("some error")
	dbConnection := &ConnectionMock{t: t, loginUserReturn: unknownError}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, authenticate)
	assert.NotNil(t, err)
	assert.Nil(t, tokenResponseDTO)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))

	assert.NotNil(t, dbConnection.loginUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.loginUserRecordArray[0].user.ID)

}

func TestLoginUser_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	unknownError := errors.New("some error")
	dbConnection := &ConnectionMock{t: t, updateRefreshTokenReturn: unknownError}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, authenticate)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))

	assert.NotNil(t, dbConnection.loginUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.loginUserRecordArray[0].user.ID)

	assert.Equal(t, userId, dbConnection.updateRefreshTokenRecordArray[0].userId)
	assert.Len(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt)

	assert.NotNil(t, err)
	assert.Nil(t, tokenResponseDTO)

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
