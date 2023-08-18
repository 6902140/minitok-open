package usal

import (
	"errors"
	"minitok/log"
	"minitok/response"

	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

var (
	Secret = []byte("TikTok")
	// TokenExpireDuration = time.Hour * 2 过期时间
)

// 定义JWT声明的结构体
type JWTClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

//UserId：表示用户的唯一标识符，类型为 int64。
//Username：表示用户名，类型为 string。
//jwt.RegisteredClaims：这是来自 github.com/golang-jwt/jwt 包的 RegisteredClaims 结构体，它包含了 JWT 标准规范中的一些预定义声明字段，例如 aud（接收方），exp（过期时间）等。

// 生成token
func GenToken(userid int64, userName string) (string, error) {
	claims := JWTClaims{
		UserId:   userid,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "server",
			//ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),可用于设定token过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("TikTok"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// 解析token
func ParsenToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 验证token
func VerifyToken(tokenString string) (int64, error) {

	log.Debugf("tokenString:%v", tokenString)

	if tokenString == "" {
		return int64(0), nil
	}
	claims, err := ParsenToken(tokenString)
	if err != nil {
		return int64(0), err
	}
	return claims.UserId, nil
}

//=============================gin的中间件，就是一个函数，返回gin 的HandlerFunc======================================================

func AuthMiddleware() gin.HandlerFunc {
	//生成token 并且转移控制权到下一个处理函数
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		userId, err := VerifyToken(tokenString) //验证令牌
		if err != nil || userId == int64(0) {
			response.Fail(c, "auth error", nil)
			c.Abort()
		}

		c.Set("UserId", userId)
		c.Next()
	}
}

// 部分接口不需要用户登录也可访问，如feed，pushlishlist，favList，follow/follower list
func AuthWithOutMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		tokstr := c.Query("token")

		userId, err := VerifyToken(tokstr)
		if err != nil {
			response.Fail(c, "auth error", nil)
			c.Abort()
		}

		c.Set("UserId", userId)
		c.Next() //这种情况下，context.Next() 的作用是将控制权从中间件函数传递给请求处理函数，使其能够执行其特定的功能并返回响应。这种方式允许在中间件函数中执行一些公共操作（如身份验证），然后将请求传递给实际的请求处理函数进行进一步处理。
	}
}
