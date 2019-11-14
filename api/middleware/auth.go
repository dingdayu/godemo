package middleware

import (
	"net/http"
	"strings"

	"demo/api/controller"

	"github.com/gin-gonic/gin"
)

// 处理跨域请求,支持options访问
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		token := strings.Split(authorization, " ")
		if token[0] != "Bearer" || len(token) != 2 {
			c.AbortWithStatusJSON(http.StatusForbidden, controller.ErrAuthForbidden)
			return
		}
		// 检查 token

		// 处理请求
		c.Next()
	}
}
