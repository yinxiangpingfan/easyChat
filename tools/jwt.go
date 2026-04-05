package tools

import (
	"easyChat/global"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Uuid      string `json:"uuid"`
	Telephone string `json:"telephone"`
	jwt.RegisteredClaims
}

// 生成访问令牌
func GenerateAccessToken(uuid, tel string) (string, string, error) {
	accesssClaims := MyCustomClaims{
		Uuid:      uuid,
		Telephone: tel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "easyChat",
		},
	}
	accesssToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accesssClaims).SignedString([]byte(global.Config.JWT.AccessSecret))
	if err != nil {
		return "", "", err
	}
	refreshClaims := MyCustomClaims{
		Uuid:      uuid,
		Telephone: tel,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
			Issuer:    "easyChat",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(global.Config.JWT.RefreshSecret))
	if err != nil {
		return "", "", err
	}
	return accesssToken, refreshToken, nil
}

// 验证访问令牌
func VerifyAccessToken(tokenString string) ([]string, int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.AccessSecret), nil
	})
	if err != nil {
		return nil, 1, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return []string{claims.Uuid, claims.Telephone}, 0, nil
	}
	return nil, 2, errors.New("token is not valid")
}
