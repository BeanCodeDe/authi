package auth

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BeanCodeDe/SpaceLight-Auth/internal/config"
	"github.com/BeanCodeDe/SpaceLight-AuthMiddleware/authAdapter"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var (
	signKey *rsa.PrivateKey
)

func Init() {
	err := authAdapter.Init()
	if err != nil {
		log.Fatalf("Error while init authAdapter: %v", err)
	}

	signBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Error while reading private Key: %v", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatalf("Error while parsing private Key: %v", err)
	}
}

func CreateJWTToken(userId uuid.UUID, roles []string) (string, error) {
	log.Debugf("create JWT token")

	expirationTime := callcExpirationTime(roles)

	claims := &authAdapter.Claims{
		UserId: userId,
		Roles:  roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("token creation failed: %v", err)
	}

	log.Debugf("JWT token created")
	return signedToken, nil
}

func callcExpirationTime(roles []string) time.Time {
	for _, role := range roles {
		if role == authAdapter.ServiceRole {
			return time.Now().Add(24 * time.Hour)
		}
	}
	return time.Now().Add(5 * time.Minute)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(authAdapter.AuthName)
		if authHeader == "" {
			return fmt.Errorf("no auth Header found")
		}

		claims, err := authAdapter.ParseToken(authHeader)
		if err != nil {
			return fmt.Errorf("error while parsing token: %v", err)
		}

		var token string

		if time.Now().Add(1 * time.Minute).After(time.Unix(claims.ExpiresAt, 0)) {
			token, err = CreateJWTToken(claims.UserId, claims.Roles)
			if err != nil {
				return fmt.Errorf("error while creating token: %v", err)
			}
			c.Response().Header().Set(authAdapter.AuthName, token)
		} else {
			token = authHeader
		}

		c.Set(authAdapter.ClaimName, *claims)
		c.Response().Header().Set(authAdapter.AuthName, token)
		return next(c)
	}
}
