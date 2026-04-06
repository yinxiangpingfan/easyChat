package v1

import (
	"easyChat/internal/config"
	"easyChat/internal/errors"
	"easyChat/pkg/log"
	"easyChat/pkg/tools"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MiddlewareHandler interface {
	AuthHttpMiddleware() gin.HandlerFunc
	TreaceLoggerMiddleware() gin.HandlerFunc
}

type middlewareHandler struct {
	jwtConfig config.JWTConfig
	logger    *logrus.Logger
}

func NewMiddlewareHandler(logger *logrus.Logger, jwtConfig config.JWTConfig) MiddlewareHandler {
	return &middlewareHandler{
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

func (u *middlewareHandler) AuthHttpMiddleware() gin.HandlerFunc {
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
		message, errCode, err := tools.VerifyAccessToken(u.jwtConfig, tokenString)
		if err != nil {
			reqLog := log.FromContext(c.Request.Context())
			if errCode == 2 {
				reqLog.Errorf("jwt验证发生了逻辑错误 errCode: %d, err: %v", errCode, err)
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

func (u *middlewareHandler) TreaceLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String() // 生成一个随机的 traceID
		reqLogger := u.logger.WithField("trace_id", traceID)
		newCtx := log.NewContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
