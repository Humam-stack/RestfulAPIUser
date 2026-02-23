package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint, email string) (string, error) {
	secretkey := os.Getenv("JWT_SECRET")
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_HOURS"))
	if err != nil {
		exp = 24
	}

	claims := JWTClaims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(exp) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenstring, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

func ValidateJWT(tokenstring string) (*JWTClaims, error) {
	secretkey := os.Getenv("JWT_SECRET")
	//parse dan validasi token
	token, err := jwt.ParseWithClaims(tokenstring, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method!")
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		return nil, err
	}

	//extract dari claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
