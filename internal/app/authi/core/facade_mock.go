package core

import (
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/google/uuid"
)

type (
	EmptyRecord struct {
	}

	AuthenticateRecord struct {
		UserId   uuid.UUID
		Password string
		InitUser bool
	}

	RefreshTokenRecord struct {
		UserId       uuid.UUID
		RefreshToken string
	}

	DeleteUserRecord struct {
		UserId uuid.UUID
	}

	AuthenticateResponse struct {
		TokenResponse *adapter.TokenResponseDTO
		Err           error
	}
	ErrorResponse struct {
		Err error
	}

	CoreMock struct {
		CreateUserRecordArray        []*AuthenticateRecord
		LoginUserRecordArray         []*AuthenticateRecord
		RefreshTokenRecordArray      []*RefreshTokenRecord
		UpdatePasswordRecordArray    []*AuthenticateRecord
		DeleteUserRecordArray        []*DeleteUserRecord
		DeleteInitUsersRecordArray   []*EmptyRecord
		CreateUserResponseArray      []*ErrorResponse
		LoginUserResponseArray       []*AuthenticateResponse
		RefreshTokenResponseArray    []*AuthenticateResponse
		UpdatePasswordResponseArray  []*ErrorResponse
		DeleteUserResponseArray      []*ErrorResponse
		DeleteInitUsersResponseArray []*ErrorResponse
	}
)

func (mock *CoreMock) CreateUser(userId uuid.UUID, password string, initUser bool) error {
	record := &AuthenticateRecord{UserId: userId, Password: password, InitUser: initUser}
	mock.CreateUserRecordArray = append(mock.CreateUserRecordArray, record)
	response := mock.CreateUserResponseArray[len(mock.CreateUserRecordArray)-1]
	return response.Err
}

func (mock *CoreMock) LoginUser(userId uuid.UUID, password string) (*adapter.TokenResponseDTO, error) {
	record := &AuthenticateRecord{UserId: userId, Password: password}
	mock.LoginUserRecordArray = append(mock.LoginUserRecordArray, record)
	response := mock.LoginUserResponseArray[len(mock.LoginUserRecordArray)-1]
	return response.TokenResponse, response.Err
}

func (mock *CoreMock) RefreshToken(userId uuid.UUID, refreshToken string) (*adapter.TokenResponseDTO, error) {
	record := &RefreshTokenRecord{UserId: userId, RefreshToken: refreshToken}
	mock.RefreshTokenRecordArray = append(mock.RefreshTokenRecordArray, record)
	response := mock.RefreshTokenResponseArray[len(mock.RefreshTokenRecordArray)-1]
	return response.TokenResponse, response.Err
}

func (mock *CoreMock) UpdatePassword(userId uuid.UUID, password string) error {
	record := &AuthenticateRecord{UserId: userId, Password: password}
	mock.UpdatePasswordRecordArray = append(mock.UpdatePasswordRecordArray, record)
	response := mock.UpdatePasswordResponseArray[len(mock.UpdatePasswordRecordArray)-1]
	return response.Err
}

func (mock *CoreMock) DeleteUser(userId uuid.UUID) error {
	record := &DeleteUserRecord{UserId: userId}
	mock.DeleteUserRecordArray = append(mock.DeleteUserRecordArray, record)
	response := mock.DeleteUserResponseArray[len(mock.DeleteUserRecordArray)-1]
	return response.Err
}

func (mock *CoreMock) DeleteInitUsers() error {
	record := &EmptyRecord{}
	mock.DeleteInitUsersRecordArray = append(mock.DeleteInitUsersRecordArray, record)
	response := mock.DeleteInitUsersResponseArray[len(mock.DeleteInitUsersRecordArray)-1]
	return response.Err
}
