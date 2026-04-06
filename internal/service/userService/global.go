package userService

import (
	"context"
	"easyChat/internal/config"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/repo"

	"github.com/redis/go-redis/v9"
)

type UserService interface {
	LoginService(ctx context.Context, req req.LoginRequest) (int, string, interface{}, int)
	RefreshTokenService(ctx context.Context, req req.RefreshTokenRequest) (int, string, interface{}, int)
	RegisterService(ctx context.Context, req req.RegisterRequest) (int, string, interface{}, int)
	SmsCodeService(ctx context.Context, req req.SmsCodeRequest) (int, string, interface{}, int)
	GetUserInfoService(ctx context.Context) (int, string, interface{}, int)
}

type userService struct {
	userRepo        repo.UserRepository
	redisClient     *redis.Client
	jwtConfig       config.JWTConfig
	aliAccessConfig config.AliAccessConfig
}

func NewUserService(jwtConfigs config.JWTConfig, aliAccessConfig config.AliAccessConfig, userRepo repo.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		jwtConfig:       jwtConfigs,
		userRepo:        userRepo,
		redisClient:     redisClient,
		aliAccessConfig: aliAccessConfig,
	}
}
