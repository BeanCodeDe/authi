package db

import (
	"time"

	"github.com/google/uuid"
)

type (
	DBMock struct {
		CloseRecordArray              []*CloseRecord
		CreateUserRecordArray         []*CreateUserRecord
		UpdateRefreshTokenRecordArray []*UpdateRefreshTokenRecord
		LoginUserRecordArray          []*LoginUserRecord
		CheckRefreshTokenRecordArray  []*CheckRefreshTokenRecord
		UpdatePasswordRecordArray     []*UpdatePasswordRecord
		DeleteUserRecordArray         []*DeleteUserRecord
		CreateUserResponseArray       []*ErrorResponse
		UpdateRefreshResponseArray    []*ErrorResponse
		LoginUserResponseArray        []*ErrorResponse
		CheckRefreshResponseArray     []*ErrorResponse
		UpdatePasswordResponseArray   []*ErrorResponse
		DeleteUserResponseArray       []*ErrorResponse
	}

	ErrorResponse struct {
		Err error
	}

	CloseRecord struct {
	}

	CreateUserRecord struct {
		User *UserDB
		Hash string
	}

	UpdateRefreshTokenRecord struct {
		UserId               uuid.UUID
		RefreshToken         string
		RefreshTokenExpireAt time.Time
	}

	LoginUserRecord struct {
		User *UserDB
	}
	CheckRefreshTokenRecord struct {
		UserId       uuid.UUID
		RefreshToken string
	}
	UpdatePasswordRecord struct {
		UserId   uuid.UUID
		Password string
		Hash     string
	}

	DeleteUserRecord struct {
		UserId uuid.UUID
	}
)

func (mock *DBMock) Close() {
	closeRecord := &CloseRecord{}
	mock.CloseRecordArray = append(mock.CloseRecordArray, closeRecord)
}

func (mock *DBMock) CreateUser(user *UserDB, hash string) error {
	record := &CreateUserRecord{User: user, Hash: hash}
	mock.CreateUserRecordArray = append(mock.CreateUserRecordArray, record)
	response := mock.CreateUserResponseArray[len(mock.CreateUserResponseArray)-1]
	return response.Err
}

func (mock *DBMock) UpdateRefreshToken(userId uuid.UUID, refreshToken string, refreshTokenExpireAt time.Time) error {
	record := &UpdateRefreshTokenRecord{UserId: userId, RefreshToken: refreshToken}
	mock.UpdateRefreshTokenRecordArray = append(mock.UpdateRefreshTokenRecordArray, record)
	response := mock.UpdatePasswordResponseArray[len(mock.UpdatePasswordResponseArray)-1]
	return response.Err
}

func (mock *DBMock) LoginUser(user *UserDB) error {
	record := &LoginUserRecord{User: user}
	mock.LoginUserRecordArray = append(mock.LoginUserRecordArray, record)
	response := mock.LoginUserResponseArray[len(mock.LoginUserResponseArray)-1]
	return response.Err
}

func (mock *DBMock) CheckRefreshToken(userId uuid.UUID, refreshToken string) error {
	record := &CheckRefreshTokenRecord{UserId: userId, RefreshToken: refreshToken}
	mock.CheckRefreshTokenRecordArray = append(mock.CheckRefreshTokenRecordArray, record)
	response := mock.CheckRefreshResponseArray[len(mock.CheckRefreshResponseArray)-1]
	return response.Err
}

func (mock *DBMock) UpdatePassword(userId uuid.UUID, password string, hash string) error {
	record := &UpdatePasswordRecord{UserId: userId, Password: password}
	mock.UpdatePasswordRecordArray = append(mock.UpdatePasswordRecordArray, record)
	response := mock.UpdatePasswordResponseArray[len(mock.UpdatePasswordResponseArray)-1]
	return response.Err
}

func (mock *DBMock) DeleteUser(userId uuid.UUID) error {
	record := &DeleteUserRecord{UserId: userId}
	mock.DeleteUserRecordArray = append(mock.DeleteUserRecordArray, record)
	response := mock.DeleteUserResponseArray[len(mock.DeleteUserResponseArray)-1]
	return response.Err
}
