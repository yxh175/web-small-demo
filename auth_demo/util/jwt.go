package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("secret")
var accessTokenExpireDuration = 24 * time.Second
var refreshTokenExpireDuration = 7 * 24 * time.Hour

type JWTClaims struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

// GenerateToken 签发用户Token
func GenerateToken(id uint, username string) (accessToken, refreshToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(accessTokenExpireDuration)
	rtExpireTime := nowTime.Add(refreshTokenExpireDuration)
	claims := JWTClaims{
		ID:       id,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "mall",
		},
	}
	// 加密并获得完整的编码后的字符串token
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: rtExpireTime.Unix(),
		Issuer:    "mall",
	}).SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// ParseToken 验证用户token
func ParseToken(token string) (*JWTClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*JWTClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func ParseRefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return
	}
	refreshClaim, err := ParseToken(rToken)
	if err != nil {
		return
	}

	if refreshClaim.ExpiresAt > time.Now().Unix() {
		// refresh未过期，分发新的token
		return GenerateToken(accessClaim.ID, accessClaim.UserName)
	}
	// refresh过期
	return "", "", errors.New("身份过期")
}
