package core

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
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
	symbols []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
)

type (
	TokenCore struct {
		AccessToken      string
		ExpiresIn        int
		RefreshToken     string
		RefreshExpiresIn int
	}
)

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
	refreshTokenExpireAt := time.Now().Add(10 * time.Minute)

	claimsToken := &authadapter.Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpireAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsToken)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return nil, fmt.Errorf("token creation failed: %v", err)
	}

	refreshToken := randomString()
	if err = db.UpdateRefreshToken(userId, refreshToken, refreshTokenExpireAt); err != nil {
		return nil, fmt.Errorf("refresh token could not be saved into database: %v", err)
	}
	return &TokenCore{AccessToken: signedToken, ExpiresIn: int(tokenExpireAt), RefreshToken: refreshToken, RefreshExpiresIn: int(refreshTokenExpireAt.Unix())}, nil
}

func randomString() string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, 32)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
