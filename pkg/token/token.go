package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	Salt string
}

func NewTokenService(salt string) *TokenService {
	return &TokenService{
		Salt: salt,
	}
}

func (t *TokenService) GenerateToken(expiresAt time.Duration, userId uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expiresAt).Unix(),
		Subject:   strconv.Itoa(int(userId)),
	})

	return token.SignedString([]byte(t.Salt))
}

func (t *TokenService) ParseToken(token string) (uint64, error) {
	resToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("uexpected signing method")
		}

		return []byte(t.Salt), nil
	})
	if err != nil {
		return 0, errors.New("invalid token")
	}

	if !resToken.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := resToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token")
	}

	res := claims["sub"].(string)
	id, err := strconv.Atoi(res)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return uint64(id), nil
}
