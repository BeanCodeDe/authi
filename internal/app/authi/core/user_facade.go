package core

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/internal/app/authi/util"
	"github.com/BeanCodeDe/authi/pkg/adapter"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v3"
)

type (
	UserFacade struct {
		dbConnection           db.Connection
		signKey                *rsa.PrivateKey
		accessTokenExpireTime  int
		refreshTokenExpireTime int
	}
	initUser struct {
		Id       uuid.UUID `yaml:"id" json:"id"`
		Password string    `yaml:"password" json:"password"`
	}
	configYml struct {
		Users []*initUser `yaml:"users"`
	}
)

const (
	EnvPrivateKeyPath         = "PRIVATE_KEY_PATH"
	EnvAccessTokenExpireTime  = "ACCESS_TOKEN_EXPIRE_TIME"
	EnvRefreshTokenExpireTime = "REFRESH_TOKEN_EXPIRE_TIME"
)

func NewUserFacade() (*UserFacade, error) {
	signKey, err := loadSignKey()
	if err != nil {
		return nil, err
	}

	dbConnection, err := db.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("error while initializing database: %v", err)
	}

	accessTokenExpireTime, err := util.GetEnvIntWithFallback(EnvAccessTokenExpireTime, 5)
	if err != nil {
		return nil, fmt.Errorf("error loading access token expire time from environment: %w", err)
	}
	refreshTokenExpireTime, err := util.GetEnvIntWithFallback(EnvRefreshTokenExpireTime, 10)
	if err != nil {
		return nil, fmt.Errorf("error loading refresh token expire time from environment: %w", err)
	}
	return &UserFacade{dbConnection, signKey, accessTokenExpireTime, refreshTokenExpireTime}, nil
}

func loadSignKey() (*rsa.PrivateKey, error) {
	path := util.GetEnvWithFallback(EnvPrivateKeyPath, "/token/jwtRS256.key")

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

func (userFacade *UserFacade) initDefaultUser() error {
	initUserFile := util.GetEnvWithFallback("INIT_USER_FILE", "authi.conf")

	_, err := os.Stat(initUserFile)

	if os.IsNotExist(err) {
		log.Warnf("File %s don't exist", initUserFile)
		return nil
	}

	data, err := os.ReadFile(initUserFile)
	if err != nil {
		return fmt.Errorf("error while loading init user File [%s]: %v", initUserFile, err)
	}

	var config *configYml

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("error while parsing init user File [%s]: %v", initUserFile, err)
	}

	for _, user := range config.Users {
		if err := userFacade.CreateUser(user.Id, user.Password); err != nil {
			return fmt.Errorf("error while creating init user [%v]: %v", user.Id, err)
		}
	}
	return nil
}

func (userFacade *UserFacade) CreateUser(userId uuid.UUID, password string) error {

	creationTime := time.Now()

	dbUser := &db.UserDB{ID: userId, Password: password, CreatedOn: creationTime, LastLogin: creationTime}

	if err := userFacade.dbConnection.CreateUser(dbUser, randomString()); err != nil {
		if errors.Is(err, db.ErrUserAlreadyExists) {
			if err := userFacade.dbConnection.LoginUser(dbUser); err != nil {
				return fmt.Errorf("something went wrong while checking credentials of already created user, %v: %w", userId, err)
			}
			return nil
		}
		return fmt.Errorf("error while creating user: %v", err)
	}

	return nil
}

func (userFacade *UserFacade) LoginUser(userId uuid.UUID, password string) (*adapter.TokenResponseDTO, error) {
	dbUser := &db.UserDB{ID: userId, Password: password}
	if err := userFacade.dbConnection.LoginUser(dbUser); err != nil {
		return nil, fmt.Errorf("something went wrong when logging in user, %v: %v", userId, err)
	}
	return userFacade.createJWTToken(userId)
}

func (userFacade *UserFacade) RefreshToken(userId uuid.UUID, refreshToken string) (*adapter.TokenResponseDTO, error) {
	if err := userFacade.dbConnection.CheckRefreshToken(userId, refreshToken); err != nil {
		return nil, fmt.Errorf("no user with refresh token was found: %v", err)
	}

	return userFacade.createJWTToken(userId)
}

func (userFacade *UserFacade) UpdatePassword(userId uuid.UUID, password string) error {
	if err := userFacade.dbConnection.UpdatePassword(userId, password, randomString()); err != nil {
		return fmt.Errorf("error while updating password of user: %v", err)
	}
	return nil
}

func (userFacade *UserFacade) DeleteUser(userId uuid.UUID) error {
	if err := userFacade.dbConnection.DeleteUser(userId); err != nil {
		return fmt.Errorf("error while deleting user: %v", err)
	}
	return nil
}

func (userFacade *UserFacade) createJWTToken(userId uuid.UUID) (*adapter.TokenResponseDTO, error) {

	tokenExpireAt := time.Now().Add(time.Duration(userFacade.accessTokenExpireTime) * time.Minute).Unix()
	refreshTokenExpireAt := time.Now().Add(time.Duration(userFacade.refreshTokenExpireTime) * time.Minute)

	claimsToken := &adapter.Claims{
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
	return &adapter.TokenResponseDTO{AccessToken: signedToken, ExpiresIn: int(tokenExpireAt), RefreshToken: refreshToken, RefreshExpiresIn: int(refreshTokenExpireAt.Unix())}, nil
}
