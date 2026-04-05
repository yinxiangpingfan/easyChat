package v1

import (
	"easyChat/errors"
	"easyChat/global"
	"easyChat/tools"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthHttpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Token (标准格式为 "Authorization: Bearer <token>")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			JsonBack(c, errors.ErrAuthFormat.Code, errors.ErrAuthFormat.Message, nil, -1)
			c.Abort()
			return
		}

		// 截取 Bearer 后的内容
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			JsonBack(c, errors.ErrAuthFormat.Code, errors.ErrAuthFormat.Message, nil, -1)
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 解析 Access Token
		message, errCode, err := tools.VerifyAccessToken(tokenString)
		if err != nil {
			if errCode == 2 {
				global.Logger.Errorf("jwt验证发生了逻辑错误 errCode: %d, err: %v", errCode, err)
				JsonBack(c, errors.ErrAuthFailed.Code, errors.ErrAuthFailed.Message, nil, -1)
			}
			JsonBack(c, errors.ErrAuthExpired.Code, errors.ErrAuthExpired.Message, nil, -1)
			c.Abort()
			return
		}
		//讲信息挂载到上下文，方便后续的 Handler 使用
		c.Set("uuid", message[0])
		c.Set("tel", message[1])
		c.Next()
	}
}
