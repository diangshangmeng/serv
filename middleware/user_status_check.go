package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func UserStatusCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint64)

		user, err := repository.GetUserByID(userID)
		if err != nil {
			util.ResponseError(c, http.StatusInternalServerError, "获取用户信息失败")
			c.Abort()
			return
		}

		if user.Status == 0 {
			util.ResponseError(c, http.StatusForbidden, "账号已被禁用")
			c.Abort()
			return
		}

		c.Set("userStatus", user.Status)
		c.Set("userAuthStatus", user.AuthStatus)
		c.Next()
	}
}
