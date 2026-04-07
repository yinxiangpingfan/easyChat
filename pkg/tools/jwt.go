package tools

import (
	"easyChat/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Uuid      string `json:"uuid"`
	Telephone string `json:"telephone"`
	IsAdmin   int8   `json:"is_admin"`
	jwt.RegisteredClaims
}

// 生成访问令牌
func GenerateAccessToken(config config.JWTConfig, uuid, tel string, isAdmin int8) (string, string, error) {
	accesssClaims := MyCustomClaims{
		Uuid:      uuid,
		Telephone: tel,
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "easyChat",
		},
	}
	accesssToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accesssClaims).SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := MyCustomClaims{
		Uuid:      uuid,
		Telephone: tel,
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			Issuer:    "easyChat",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return "", "", err
	}
	return accesssToken, refreshToken, nil
}

// 验证访问令牌
func VerifyAccessToken(config config.JWTConfig, tokenString string) ([]string, int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AccessSecret), nil
	})
	if err != nil {
		return nil, 1, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		user := "user"
		switch claims.IsAdmin {
		case 1:
			user = "admin"
		case 2:
			user = "superAdmin"
		}
		return []string{claims.Uuid, claims.Telephone, user}, 0, nil
	}
	return nil, 2, errors.New("token is not valid")
}

// 验证刷新令牌
func VerifyRefreshToken(config config.JWTConfig, tokenString string) ([]string, int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshSecret), nil
	})
	if err != nil {
		return nil, 1, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		user := "user"
		switch claims.IsAdmin {
		case 1:
			user = "admin"
		case 2:
			user = "superAdmin"
		}
		return []string{claims.Uuid, claims.Telephone, user}, 0, nil
	}
	return nil, 2, errors.New("token is not valid")
}
