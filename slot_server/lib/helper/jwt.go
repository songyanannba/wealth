package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const JwtSecretKey = "gateWaySys2025"

type JWTInstance struct {
	SecretKey []byte
}

func InitJwt(SecretKey []byte) JWTInstance {
	return JWTInstance{SecretKey}
}

func (that JWTInstance) GenerateJWT(uid, sub string, count time.Duration) (string, error) {
	// 创建一个新的JWT token

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	// 设置一些声明
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["user_id"] = uid
	claims["exp"] = time.Now().Add(time.Hour * count).Unix()
	claims["sub"] = sub

	// 设置签名并获取token字符串
	token, err := jwtToken.SignedString(that.SecretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (that JWTInstance) ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// 解析JWT字符串
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return that.SecretKey, nil
	})

	if err != nil {
		return nil, errors.New("jwt过期")
	}

	// 验证token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("jwt token err")
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	jwt := InitJwt([]byte(JwtSecretKey))
	return jwt.ParseJWT(tokenString)
}

func GenerateJWT(uid, sub string, count time.Duration) (string, error) {
	jwt := InitJwt([]byte(JwtSecretKey))
	return jwt.GenerateJWT(uid, sub, count)
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) (string, error) {
	if claims, exists := c.Get("claims"); !exists {
		return "", errors.New("user_id not exist")
	} else {
		waitUse := claims.(jwt.MapClaims)
		//fmt.Println(waitUse)
		return waitUse["user_id"].(string), nil
	}
}
