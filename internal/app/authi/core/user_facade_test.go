package core

import (
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/stretchr/testify/assert"
)

// CreateUser Test

func TestCreateUser_Successfully(t *testing.T) {
	dbConnection := &ConnectionMock{}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, authenticate)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.NotNil(t, dbConnection.createUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
	assert.Equal(t, authenticate.Password, dbConnection.createUserRecordArray[0].user.Password)
	assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)

}

func TestCreateUser_CreateUser_UnknownError(t *testing.T) {
	dbConnection := &ConnectionMock{createUserReturn: errUnknown}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, authenticate)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
	assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
}

func TestCreateUser_CreateUser_AlreadyExists(t *testing.T) {
	dbConnection := &ConnectionMock{createUserReturn: errormessages.ErrUserAlreadyExists}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, authenticate)
	assert.Equal(t, errormessages.ErrUserAlreadyExists, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 1, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.createUserRecordArray[0].user.ID)
	assert.Len(t, dbConnection.createUserRecordArray[0].hash, 32)
}

// RefreshToken Test

func TestRefreshToken_Successfully(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &ConnectionMock{}

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
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

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
	dbConnection := &ConnectionMock{checkRefreshReturn: errUnknown}

	userFacade := &UserFacade{dbConnection: dbConnection}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.checkRefreshTokenRecordArray[0].userId)
	assert.Equal(t, refreshToken, dbConnection.checkRefreshTokenRecordArray[0].refreshToken)

	assert.Nil(t, tokenResponseDTO)
	assert.NotNil(t, err)
}

func TestRefreshToken_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &ConnectionMock{updateRefreshTokenReturn: errUnknown}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

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
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &ConnectionMock{}

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
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

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
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)

	dbConnection := &ConnectionMock{loginUserReturn: errUnknown}

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
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.NotNil(t, dbConnection.loginUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.loginUserRecordArray[0].user.ID)

}

func TestLoginUser_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &ConnectionMock{updateRefreshTokenReturn: errUnknown}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, authenticate)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.NotNil(t, dbConnection.loginUserRecordArray[0].user)
	assert.Equal(t, userId, dbConnection.loginUserRecordArray[0].user.ID)

	assert.Equal(t, userId, dbConnection.updateRefreshTokenRecordArray[0].userId)
	assert.Len(t, dbConnection.updateRefreshTokenRecordArray[0].refreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.updateRefreshTokenRecordArray[0].refreshTokenExpireAt)

	assert.NotNil(t, err)
	assert.Nil(t, tokenResponseDTO)

}

// UpdatePassword Test
func TestUpdatePassword_Successfully(t *testing.T) {
	dbConnection := &ConnectionMock{}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.UpdatePassword(userId, authenticate)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.updatePasswordRecordArray[0].userId)
	assert.Equal(t, authenticate.Password, dbConnection.updatePasswordRecordArray[0].password)
	assert.Len(t, dbConnection.updatePasswordRecordArray[0].hash, 32)
}

func TestUpdatePassword_UnknownError(t *testing.T) {
	dbConnection := &ConnectionMock{updatePasswordReturn: errUnknown}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.UpdatePassword(userId, authenticate)

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.updatePasswordRecordArray[0].userId)
	assert.Equal(t, authenticate.Password, dbConnection.updatePasswordRecordArray[0].password)
	assert.Len(t, dbConnection.updatePasswordRecordArray[0].hash, 32)
}

// DeleteUser Test

func TestDeleteUser_Successfully(t *testing.T) {
	dbConnection := &ConnectionMock{}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.DeleteUser(userId)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 1, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.deleteUserRecordArray[0].userId)
}

func TestDeleteUser_DeleteUser_UnknownError(t *testing.T) {
	dbConnection := &ConnectionMock{deleteUserReturn: errUnknown}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.DeleteUser(userId)

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.closeRecordArray))
	assert.Equal(t, 0, len(dbConnection.createUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.updateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.loginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.checkRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.updatePasswordRecordArray))
	assert.Equal(t, 1, len(dbConnection.deleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.deleteUserRecordArray[0].userId)
}
