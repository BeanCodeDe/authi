package core

import (
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type (
	ConnectionMock struct {
		t                             *testing.T
		closeRecordArray              []*closeRecord
		createUserRecordArray         []*createUserRecord
		updateRefreshTokenRecordArray []*updateRefreshTokenRecord
		getUserByIdRecordArray        []*getUserByIdRecord
		loginUserRecordArray          []*loginUserRecord
		checkRefreshTokenRecordArray  []*checkRefreshTokenRecord
	}
	AuthMock struct {
		t                         *testing.T
		parseTokenRecordArray     []*parseTokenRecord
		createJWTTokenRecordArray []*createJWTTokenRecord
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

	getUserByIdRecord struct {
		userId uuid.UUID
	}
	loginUserRecord struct {
		user *db.UserDB
	}
	checkRefreshTokenRecord struct {
		userId       uuid.UUID
		refreshToken string
	}

	parseTokenRecord struct {
		authorizationString string
	}

	createJWTTokenRecord struct {
		token string
	}
)

var (
	userId         = uuid.New()
	refreshToken   = "someRefreshToken"
	authenticate   = &AuthenticateDTO{}
	privateKeyPath = "../../../../deployments/token/privat/jwtRS256.key"
)

func (connection *ConnectionMock) Close() {
	closeRecord := &closeRecord{}
	connection.closeRecordArray = append(connection.closeRecordArray, closeRecord)
}

// CreateUser Test

func TestCreateUser_Successfully(t *testing.T) {
	dbConnection := &ConnectionMock{t: t}
	auth := &AuthMock{t: t}
	userFacade := &UserFacade{dbConnection: dbConnection, authAdapter: auth}

	if assert.NoError(t, userFacade.CreateUser(userId, authenticate)) {
		assert.Equal(t, 0, len(dbConnection.closeRecordArray))
		assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
		assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
		assert.Equal(t, 0, len(dbConnection.getUserByIdRecordArray))
		assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
		assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
		assert.Equal(t, 0, len(auth.parseTokenRecordArray))
		assert.Equal(t, 0, len(auth.createJWTTokenRecordArray))

		assert.NotNil(t, dbConnection.createUserRecordArray[0].user)
		assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
		assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
	}
}

func (connection *ConnectionMock) CreateUser(user *db.UserDB, hash string) error {
	createUserRecord := &createUserRecord{user: user, hash: hash}
	connection.createUserRecordArray = append(connection.createUserRecordArray, createUserRecord)
	return nil
}

// Test refresh token

func TestRefreshToken_Successfully(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	dbConnection := &ConnectionMock{t: t}
	auth := &AuthMock{t: t}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, authAdapter: auth, signKey: signKey}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, len(dbConnection.closeRecordArray))
		assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
		assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
		assert.Equal(t, 0, len(dbConnection.getUserByIdRecordArray))
		assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
		assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))
		assert.Equal(t, 0, len(auth.parseTokenRecordArray))
		assert.Equal(t, 0, len(auth.createJWTTokenRecordArray))

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
}

func (connection *ConnectionMock) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	updateRefreshToken := &updateRefreshTokenRecord{userId: userId, refreshToken: refreshToken, refreshTokenExpireAt: refreshTokenExpireAt}
	connection.updateRefreshTokenRecordArray = append(connection.updateRefreshTokenRecordArray, updateRefreshToken)
	return nil
}
func (connection *ConnectionMock) GetUserById(userId uuid.UUID) (*db.UserDB, error) {
	getUserByIdRecord := &getUserByIdRecord{userId: userId}
	connection.getUserByIdRecordArray = append(connection.getUserByIdRecordArray, getUserByIdRecord)
	return &db.UserDB{}, nil
}

// LoginUser Test

func TestLoginUser_Successfully(t *testing.T) {
	t.Setenv(PRIVATE_KEY_PATH_ENV, privateKeyPath)
	dbConnection := &ConnectionMock{t: t}
	auth := &AuthMock{t: t}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, authAdapter: auth, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, authenticate)
	if assert.Nil(t, err) {
		assert.Equal(t, 0, len(dbConnection.closeRecordArray))
		assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
		assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
		assert.Equal(t, 0, len(dbConnection.getUserByIdRecordArray))
		assert.Equal(t, 1, len(dbConnection.loginUserRecordArray))
		assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
		assert.Equal(t, 0, len(auth.parseTokenRecordArray))
		assert.Equal(t, 0, len(auth.createJWTTokenRecordArray))

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
}

func (connection *ConnectionMock) LoginUser(user *db.UserDB) error {
	loginUserRecord := &loginUserRecord{user: user}
	connection.loginUserRecordArray = append(connection.loginUserRecordArray, loginUserRecord)
	return nil
}
func (connection *ConnectionMock) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {
	checkRefreshTokenRecord := &checkRefreshTokenRecord{userId: userId, refreshToken: refreshToken}
	connection.checkRefreshTokenRecordArray = append(connection.checkRefreshTokenRecordArray, checkRefreshTokenRecord)
	return nil
}

func (auth *AuthMock) ParseToken(authorizationString string) (*authadapter.Claims, error) {
	parseTokenRecord := &parseTokenRecord{authorizationString: authorizationString}
	auth.parseTokenRecordArray = append(auth.parseTokenRecordArray, parseTokenRecord)
	return &authadapter.Claims{}, nil
}
func (auth *AuthMock) CreateJWTToken(token string) (string, error) {
	createJWTTokenRecord := &createJWTTokenRecord{token: token}
	auth.createJWTTokenRecordArray = append(auth.createJWTTokenRecordArray, createJWTTokenRecord)
	return "", nil
}
