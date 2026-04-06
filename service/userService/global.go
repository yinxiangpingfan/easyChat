package userService

import (
	"context"
	"easyChat/handler/v1/req"
	"easyChat/repo"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserService interface {
	LoginService(c *gin.Context, req req.LoginRequest) (int, string, interface{}, int)
	RefreshTokenService(c *gin.Context, req req.RefreshTokenRequest) (int, string, interface{}, int)
	RegisterService(ctx context.Context, req req.RegisterRequest) (int, string, interface{}, int)
	SmsCodeService(ctx context.Context, req req.SmsCodeRequest) (int, string, interface{}, int)
}

type userService struct {
	userRepo    repo.UserRepository
	redisClient *redis.Client
}

func NewUserService(userRepo repo.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}
