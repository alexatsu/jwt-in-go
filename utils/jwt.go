package utils

import (
	"fmt"
	"os"
	"time"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	token          *jwt.Token
	expirationTime time.Time
	secret         = []byte(os.Getenv("JWT_SECRET"))
)

// move to const time variables
var (
	TimeFiveMinutes = time.Now().Add(5 * time.Minute)
	TimeOneMonth    = time.Now().Add(30 * 24 * time.Hour)
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Exp    string      `json:"exp"`
	UserId pgtype.UUID `json:"userId"`
}

func genJWT(tokenType string, userId pgtype.UUID) (string, error) {
	if tokenType == "access" {
		expirationTime = TimeFiveMinutes
	} else if tokenType == "refresh" {
		expirationTime = TimeOneMonth
	}

	token = jwt.New(jwt.SigningMethodHS256)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Exp:    expirationTime.Format(time.RFC3339),
		UserId: userId,
	})

	str, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return str, err
}

func ValidateToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		return nil, errors.New("invalid type for claims")
	}

	return claims, nil
}

func IsTokenExpired(claims *CustomClaims) bool {
	expirationTime, err := time.Parse(time.RFC3339, claims.Exp)

	if err != nil {
		return true
	}

	return expirationTime.Before(time.Now())
}

func GenerateTokens(userId pgtype.UUID) (accessToken string, refreshToken string, err error) {
	accessToken, err = genJWT("access", userId)
	if err != nil {
		return
	}

	refreshToken, err = genJWT("refresh", userId)
	if err != nil {
		return
	}

	return
}
