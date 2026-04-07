package v1

import (
	"easyChat/internal/config"
	"easyChat/internal/errors"
	"easyChat/internal/service/adminServer"
	"easyChat/pkg/log"
	"easyChat/pkg/tools"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type MiddlewareHandler interface {
	AuthHttpMiddleware() gin.HandlerFunc
	TreaceLoggerMiddleware() gin.HandlerFunc
	CheckAdminMiddleware() gin.HandlerFunc
}

type middlewareHandler struct {
	jwtConfig   config.JWTConfig
	logger      *logrus.Logger
	redisClient *redis.Client
	cache       *cache.Cache
}

func NewMiddlewareHandler(logger *logrus.Logger, jwtConfig config.JWTConfig, redisClient *redis.Client, cache *cache.Cache) MiddlewareHandler {
	return &middlewareHandler{
		jwtConfig:   jwtConfig,
		logger:      logger,
		redisClient: redisClient,
		cache:       cache,
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
		//解析到用户数据，查看用户是否被封禁
		//1.从cache中查看用户是否被封禁
		if status, ok := u.cache.Get(message[0]); ok {
			if status.(string) == "banned" {
				reqLog := log.FromContext(c.Request.Context())
				reqLog.Infof("用户 %s 被封禁", message[0])
				JsonBack(c, errors.ErrLoginUserBanned.Code, errors.ErrLoginUserBanned.Message, nil, -1)
				c.Abort()
				return
			}
			if status.(string) == "normal" {
				//讲信息挂载到上下文，方便后续的 Handler 使用
				c.Set("uuid", message[0])
				c.Set("tel", message[1])
				c.Set("user", message[2])
				c.Next()
				return
			}
		}
		//2.从redis中查看用户是否被封禁
		status, err := u.redisClient.Exists(c.Request.Context(), adminServer.BannedKey+message[0]).Result()
		if err != nil {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Errorf("从redis中查看用户是否被封禁失败, uuid: %s, err: %v", message[0], err)
			JsonBack(c, errors.ErrAuthFailed.Code, errors.ErrAuthFailed.Message, nil, -1)
			c.Abort()
			return
		}
		if status > 0 {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Errorf("用户 %s 被封禁", message[0])
			//将用户状态缓存到cache中
			u.cache.Set(message[0], "banned", 15*time.Minute)
			JsonBack(c, errors.ErrLoginUserBanned.Code, errors.ErrLoginUserBanned.Message, nil, -1)
			c.Abort()
			return
		}
		//3.用户正常，将用户状态缓存到cache中
		c.Set("uuid", message[0])
		c.Set("tel", message[1])
		c.Set("user", message[2])
		u.cache.Set(message[0], "normal", 1*time.Minute)
		c.Next()
	}
}

// 检查用户是否是管理员
func (u *middlewareHandler) CheckAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("user")
		if !exists {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Errorf("用户不是管理员")
			c.Abort()
			return
		}
		if isAdmin.(string) != "admin" && isAdmin.(string) != "superAdmin" {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Errorf("用户不是管理员")
			c.Abort()
			return
		}
		c.Next()
	}
}

// TreaceLoggerMiddleware 记录请求日志	并挂载到上下文
func (u *middlewareHandler) TreaceLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String() // 生成一个随机的 traceID
		reqLogger := u.logger.WithField("trace_id", traceID)
		newCtx := log.NewContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
