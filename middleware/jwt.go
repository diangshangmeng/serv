package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"voucher-platform/config"
	"voucher-platform/util"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.ResponseError(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.ResponseError(c, http.StatusUnauthorized, "Token格式错误")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			util.ResponseError(c, http.StatusUnauthorized, "Token无效")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.ResponseError(c, http.StatusUnauthorized, "Token解析失败")
			c.Abort()
			return
		}

		userIDFloat, _ := claims["user_id"].(float64)
		userID := uint64(userIDFloat)
		c.Set("userID", userID)
		c.Set("phone", claims["phone"])
		c.Next()
	}
}

func AdminJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.ResponseError(c, http.StatusUnauthorized, "未登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.ResponseError(c, http.StatusUnauthorized, "Token格式错误")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			util.ResponseError(c, http.StatusUnauthorized, "Token无效")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.ResponseError(c, http.StatusUnauthorized, "Token解析失败")
			c.Abort()
			return
		}

		if claims["role"] != "admin" {
			util.ResponseError(c, http.StatusForbidden, "无权限")
			c.Abort()
			return
		}

		adminIDFloat, _ := claims["admin_id"].(float64)
		adminID := uint64(adminIDFloat)
		c.Set("adminID", adminID)
		c.Set("adminPhone", claims["username"])
		c.Next()
	}
}