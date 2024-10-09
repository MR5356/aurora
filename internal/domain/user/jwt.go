package user

import (
	"errors"
	"github.com/MR5356/aurora/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	onceJWT    sync.Once
	jwtService *JWTService
)

type JWTService struct {
	jwt config.JWT
}

func NewJWTService(cfg *config.Config) *JWTService {
	onceJWT.Do(func() {
		jwtService = &JWTService{
			jwt: cfg.JWT,
		}
	})
	return jwtService
}

func GetJWTService() *JWTService {
	return jwtService
}

func (s *JWTService) CreateToken(user *User) (string, error) {
	logrus.Debugf("create token for user: %+v", user)
	now := time.Now()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    s.jwt.Issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(s.jwt.Expire)),
				NotBefore: jwt.NewNumericDate(now),
				ID:        user.ID,
			},
		},
	)
	return token.SignedString([]byte(s.jwt.Secret))
}

func (s *JWTService) ParseToken(tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.User, nil
	}
	return nil, errors.New("invalid token")
}

type CustomClaims struct {
	User *User
	jwt.RegisteredClaims
}
