package core

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/BeanCodeDe/authi/internal/app/authi/config"
	"github.com/BeanCodeDe/authi/internal/app/authi/db"
	"github.com/BeanCodeDe/authi/pkg/authadapter"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	signKey *rsa.PrivateKey
)

type TokenCore struct {
	AccessToken      string
	ExpiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
}

func Init() {
	err := authadapter.Init()
	if err != nil {
		log.Fatalf("Error while init authAdapter: %v", err)
	}

	signBytes, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Error while reading private Key: %v", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error while parsing private Key: %v", err)
	}
}

func CreateJWTTokenFromRefreshToken(userId uuid.UUID, refreshToken string) (*TokenCore, error) {
	if err := db.CheckRefreshToken(userId, refreshToken); err != nil {
		return nil, fmt.Errorf("no user with refreshtoken was found: %v", err)
	}

	return createJWTToken(userId)
}

func createJWTToken(userId uuid.UUID) (*TokenCore, error) {

	tokenExpireAt := time.Now().Add(5 * time.Minute).Unix()
	refreshTokenExpireAt := time.Now().Add(10 * time.Minute).Unix()

	claimsToken := &authadapter.Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpireAt,
		},
	}

	claimsRefreshToken := &authadapter.Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpireAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsToken)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsRefreshToken)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return nil, fmt.Errorf("token creation failed: %v", err)
	}

	signedRefreshToken, err := refreshToken.SignedString(signKey)
	if err != nil {
		return nil, fmt.Errorf("refresh token creation failed: %v", err)
	}

	if err = db.UpdateRefreshToken(userId, signedRefreshToken); err != nil {
		return nil, fmt.Errorf("refresh token could not be saved into database: %v", err)
	}
	return &TokenCore{AccessToken: signedToken, ExpiresIn: int(tokenExpireAt), RefreshToken: signedRefreshToken, RefreshExpiresIn: int(refreshTokenExpireAt)}, nil
}
