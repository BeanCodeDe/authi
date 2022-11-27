package core

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/internal/app/authi/errormessages"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	UserFacade struct {
		dbConnection db.Connection
		authAdapter  authadapter.Auth
		signKey      *rsa.PrivateKey
	}
)

const (
	PRIVATE_KEY_PATH_ENV = "PRIVATE_KEY_PATH"
)

func NewUserFacade(authAdapter authadapter.Auth) (*UserFacade, error) {
	signKey, err := loadSignKey()
	if err != nil {
		return nil, err
	}

	dbConnection, err := db.NewPostgresConnection()
	if err != nil {
		return nil, fmt.Errorf("error while initializing database: %v", err)
	}
	return &UserFacade{dbConnection, authAdapter, signKey}, nil
}

func loadSignKey() (*rsa.PrivateKey, error) {
	path := os.Getenv(PRIVATE_KEY_PATH_ENV)
	signBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading private Key: %v", err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, fmt.Errorf("error while parsing private Key: %v", err)
	}

	return signKey, nil
}

func (userFacade *UserFacade) CreateUser(userId uuid.UUID, authenticate *AuthenticateDTO) error {

	creationTime := time.Now()

	dbUser := &db.UserDB{ID: userId, Password: authenticate.Password, CreatedOn: creationTime, LastLogin: creationTime}

	if err := userFacade.dbConnection.CreateUser(dbUser, randomString()); err != nil {
		if errors.Is(err, errormessages.ErrUserAlreadyExists) {
			return err
		}
		return fmt.Errorf("error while creating user: %v", err)
	}

	return nil
}

func (userFacade *UserFacade) LoginUser(userId uuid.UUID, authenticate *AuthenticateDTO) (*TokenResponseDTO, error) {
	dbUser := &db.UserDB{ID: userId, Password: authenticate.Password}
	if err := userFacade.dbConnection.LoginUser(dbUser); err != nil {
		return nil, fmt.Errorf("something went wrong when logging in user, %v: %v", userId, err)
	}
	return userFacade.createJWTToken(userId)
}

func (userFacade *UserFacade) RefreshToken(userId uuid.UUID, refreshToken string) (*TokenResponseDTO, error) {
	if err := userFacade.dbConnection.CheckRefreshToken(userId, refreshToken); err != nil {
		return nil, fmt.Errorf("no user with refresh token was found: %v", err)
	}

	return userFacade.createJWTToken(userId)
}

func (userFacade *UserFacade) createJWTToken(userId uuid.UUID) (*TokenResponseDTO, error) {

	tokenExpireAt := time.Now().Add(5 * time.Minute).Unix()
	refreshTokenExpireAt := time.Now().Add(10 * time.Minute)

	claimsToken := &authadapter.Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpireAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsToken)
	signedToken, err := token.SignedString(userFacade.signKey)
	if err != nil {
		return nil, fmt.Errorf("token creation failed: %v", err)
	}

	refreshToken := randomString()
	if err = userFacade.dbConnection.UpdateRefreshToken(userId, refreshToken, refreshTokenExpireAt); err != nil {
		return nil, fmt.Errorf("refresh token could not be saved into database: %v", err)
	}
	return &TokenResponseDTO{AccessToken: signedToken, ExpiresIn: int(tokenExpireAt), RefreshToken: refreshToken, RefreshExpiresIn: int(refreshTokenExpireAt.Unix())}, nil
}
