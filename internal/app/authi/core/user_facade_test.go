package core

import (
	"os"
	"testing"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/stretchr/testify/assert"
)

// InitUser Test
func TestInitUser_YamlSuccessfully(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: nil}, {Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	os.Setenv(EnvInitUserFile, "user_test.yml")
	err := userFacade.initDefaultUser()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 2, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.CreateUserRecordArray[0].User)
	assert.Equal(t, "c5ffc340-507e-4c66-a6ce-a7d98842f9ba", dbConnection.CreateUserRecordArray[0].User.ID.String())
	assert.Equal(t, "someSecretPassword", dbConnection.CreateUserRecordArray[0].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)

	assert.NotNil(t, dbConnection.CreateUserRecordArray[1].User)
	assert.Equal(t, "5cc3621d-e5ac-4d81-93df-462b27e0cc2b", dbConnection.CreateUserRecordArray[1].User.ID.String())
	assert.Equal(t, "someOtherPassword", dbConnection.CreateUserRecordArray[1].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[1].Hash, 32)
}

func TestInitUser_JsonSuccessfully(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: nil}, {Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	os.Setenv(EnvInitUserFile, "user_test.json")
	err := userFacade.initDefaultUser()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 2, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.CreateUserRecordArray[0].User)
	assert.Equal(t, "c5ffc340-507e-4c66-a6ce-a7d98842f9ba", dbConnection.CreateUserRecordArray[0].User.ID.String())
	assert.Equal(t, "someSecretPassword", dbConnection.CreateUserRecordArray[0].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)

	assert.NotNil(t, dbConnection.CreateUserRecordArray[1].User)
	assert.Equal(t, "5cc3621d-e5ac-4d81-93df-462b27e0cc2b", dbConnection.CreateUserRecordArray[1].User.ID.String())
	assert.Equal(t, "someOtherPassword", dbConnection.CreateUserRecordArray[1].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[1].Hash, 32)
}

func TestInitUser_FileNotFound(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: nil}, {Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	os.Setenv(EnvInitUserFile, "some_random_file.yml")
	err := userFacade.initDefaultUser()

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))
}

func TestInitUser_NotPasable(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: nil}, {Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	os.Setenv(EnvInitUserFile, "user_test_wrong_fornat.yml")
	err := userFacade.initDefaultUser()

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))
}

// CreateUser Test

func TestCreateUser_Successfully(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, password, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 1, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.CreateUserRecordArray[0].User)
	assert.Equal(t, userId, dbConnection.CreateUserRecordArray[0].User.ID)
	assert.Equal(t, false, dbConnection.CreateUserRecordArray[0].User.InitUser)
	assert.Equal(t, authenticate.Password, dbConnection.CreateUserRecordArray[0].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)

}

func TestCreateUser_CreateUser_UnknownError(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, password, false)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 1, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CreateUserRecordArray[0].User.ID)
	assert.Equal(t, false, dbConnection.CreateUserRecordArray[0].User.InitUser)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)
}

func TestCreateUser_CreateUser_AlreadyExistsWrongPassword(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: db.ErrUserAlreadyExists}}, LoginUserResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, password, false)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 1, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CreateUserRecordArray[0].User.ID)
	assert.Equal(t, false, dbConnection.CreateUserRecordArray[0].User.InitUser)
	assert.Equal(t, password, dbConnection.CreateUserRecordArray[0].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)

	assert.Equal(t, userId, dbConnection.LoginUserRecordArray[0].User.ID)
	assert.Equal(t, password, dbConnection.LoginUserRecordArray[0].User.Password)
}

func TestCreateUser_CreateUser_AlreadyExistsRetry(t *testing.T) {
	dbConnection := &db.DBMock{CreateUserResponseArray: []*db.ErrorResponse{{Err: db.ErrUserAlreadyExists}}, LoginUserResponseArray: []*db.ErrorResponse{{Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.CreateUser(userId, password, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 1, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CreateUserRecordArray[0].User.ID)
	assert.Equal(t, false, dbConnection.CreateUserRecordArray[0].User.InitUser)
	assert.Equal(t, password, dbConnection.CreateUserRecordArray[0].User.Password)
	assert.Len(t, dbConnection.CreateUserRecordArray[0].Hash, 32)

	assert.Equal(t, userId, dbConnection.LoginUserRecordArray[0].User.ID)
	assert.Equal(t, password, dbConnection.LoginUserRecordArray[0].User.Password)
}

// RefreshToken Test

func TestRefreshToken_Successfully(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &db.DBMock{UpdateRefreshTokenResponseArray: []*db.ErrorResponse{{Err: nil}}, CheckRefreshTokenResponseArray: []*db.ErrorResponse{{Err: nil}}}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey, accessTokenExpireTime: 5, refreshTokenExpireTime: 10}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CheckRefreshTokenRecordArray[0].UserId)
	assert.Equal(t, refreshToken, dbConnection.CheckRefreshTokenRecordArray[0].RefreshToken)

	assert.Equal(t, userId, dbConnection.UpdateRefreshTokenRecordArray[0].UserId)
	assert.Len(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt)

	assert.Regexp(t, "([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]*)", tokenResponseDTO.AccessToken)
	assert.Less(t, int(time.Now().Unix()), tokenResponseDTO.ExpiresIn)
	assert.Equal(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, tokenResponseDTO.RefreshToken)
	assert.Equal(t, int(dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt.Unix()), tokenResponseDTO.RefreshExpiresIn)
}

func TestRefreshToken_CheckRefreshToken_UnknownError(t *testing.T) {
	dbConnection := &db.DBMock{CheckRefreshTokenResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}

	userFacade := &UserFacade{dbConnection: dbConnection}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CheckRefreshTokenRecordArray[0].UserId)
	assert.Equal(t, refreshToken, dbConnection.CheckRefreshTokenRecordArray[0].RefreshToken)

	assert.Nil(t, tokenResponseDTO)
	assert.NotNil(t, err)
}

func TestRefreshToken_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &db.DBMock{UpdateRefreshTokenResponseArray: []*db.ErrorResponse{{Err: errUnknown}}, CheckRefreshTokenResponseArray: []*db.ErrorResponse{{Err: nil}}}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey, accessTokenExpireTime: 5, refreshTokenExpireTime: 10}

	tokenResponseDTO, err := userFacade.RefreshToken(userId, refreshToken)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.CheckRefreshTokenRecordArray[0].UserId)
	assert.Equal(t, refreshToken, dbConnection.CheckRefreshTokenRecordArray[0].RefreshToken)

	assert.Equal(t, userId, dbConnection.UpdateRefreshTokenRecordArray[0].UserId)
	assert.Len(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt)

	assert.Nil(t, tokenResponseDTO)
	assert.NotNil(t, err)
}

// LoginUser Test

func TestLoginUser_Successfully(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &db.DBMock{UpdateRefreshTokenResponseArray: []*db.ErrorResponse{{Err: nil}}, LoginUserResponseArray: []*db.ErrorResponse{{Err: nil}}}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey, accessTokenExpireTime: 5, refreshTokenExpireTime: 10}

	tokenResponseDTO, err := userFacade.LoginUser(userId, password)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.LoginUserRecordArray[0].User)
	assert.Equal(t, userId, dbConnection.LoginUserRecordArray[0].User.ID)

	assert.Equal(t, userId, dbConnection.UpdateRefreshTokenRecordArray[0].UserId)
	assert.Len(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt)

	assert.Regexp(t, "([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]+).([a-zA-Z0-9_=]*)", tokenResponseDTO.AccessToken)
	assert.Less(t, int(time.Now().Unix()), tokenResponseDTO.ExpiresIn)
	assert.Equal(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, tokenResponseDTO.RefreshToken)
	assert.Equal(t, int(dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt.Unix()), tokenResponseDTO.RefreshExpiresIn)

}

func TestLoginUser_LoginUser_UnknownError(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)

	dbConnection := &db.DBMock{LoginUserResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey}

	tokenResponseDTO, err := userFacade.LoginUser(userId, password)
	assert.NotNil(t, err)
	assert.Nil(t, tokenResponseDTO)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.LoginUserRecordArray[0].User)
	assert.Equal(t, userId, dbConnection.LoginUserRecordArray[0].User.ID)

}

func TestLoginUser_UpdateRefreshToken_UnknownError(t *testing.T) {
	t.Setenv(EnvPrivateKeyPath, privateKeyPath)
	dbConnection := &db.DBMock{UpdateRefreshTokenResponseArray: []*db.ErrorResponse{{Err: errUnknown}}, LoginUserResponseArray: []*db.ErrorResponse{{Err: nil}}}

	signKey, err := loadSignKey()
	assert.Nil(t, err)
	userFacade := &UserFacade{dbConnection: dbConnection, signKey: signKey, accessTokenExpireTime: 5, refreshTokenExpireTime: 10}

	tokenResponseDTO, err := userFacade.LoginUser(userId, password)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.NotNil(t, dbConnection.LoginUserRecordArray[0].User)
	assert.Equal(t, userId, dbConnection.LoginUserRecordArray[0].User.ID)

	assert.Equal(t, userId, dbConnection.UpdateRefreshTokenRecordArray[0].UserId)
	assert.Len(t, dbConnection.UpdateRefreshTokenRecordArray[0].RefreshToken, 32)
	assert.Less(t, time.Now(), dbConnection.UpdateRefreshTokenRecordArray[0].RefreshTokenExpireAt)

	assert.NotNil(t, err)
	assert.Nil(t, tokenResponseDTO)

}

// UpdatePassword Test
func TestUpdatePassword_Successfully(t *testing.T) {
	dbConnection := &db.DBMock{UpdatePasswordResponseArray: []*db.ErrorResponse{{Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.UpdatePassword(userId, password)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.UpdatePasswordRecordArray[0].UserId)
	assert.Equal(t, authenticate.Password, dbConnection.UpdatePasswordRecordArray[0].Password)
	assert.Len(t, dbConnection.UpdatePasswordRecordArray[0].Hash, 32)
}

func TestUpdatePassword_UnknownError(t *testing.T) {
	dbConnection := &db.DBMock{UpdatePasswordResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.UpdatePassword(userId, password)

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 1, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 0, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.UpdatePasswordRecordArray[0].UserId)
	assert.Equal(t, authenticate.Password, dbConnection.UpdatePasswordRecordArray[0].Password)
	assert.Len(t, dbConnection.UpdatePasswordRecordArray[0].Hash, 32)
}

// DeleteUser Test

func TestDeleteUser_Successfully(t *testing.T) {
	dbConnection := &db.DBMock{DeleteUserResponseArray: []*db.ErrorResponse{{Err: nil}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.DeleteUser(userId)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 1, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.DeleteUserRecordArray[0].UserId)
}

func TestDeleteUser_DeleteUser_UnknownError(t *testing.T) {
	dbConnection := &db.DBMock{DeleteUserResponseArray: []*db.ErrorResponse{{Err: errUnknown}}}
	userFacade := &UserFacade{dbConnection: dbConnection}

	err := userFacade.DeleteUser(userId)

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(dbConnection.CloseRecordArray))
	assert.Equal(t, 0, len(dbConnection.CreateUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdateRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.LoginUserRecordArray))
	assert.Equal(t, 0, len(dbConnection.CheckRefreshTokenRecordArray))
	assert.Equal(t, 0, len(dbConnection.UpdatePasswordRecordArray))
	assert.Equal(t, 1, len(dbConnection.DeleteUserRecordArray))

	assert.Equal(t, userId, dbConnection.DeleteUserRecordArray[0].UserId)
}
