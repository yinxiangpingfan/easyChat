package adminServer

import (
	"context"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/repo"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type AdminServer interface {
	GetUserInfoList(ctx context.Context, req *req.UserInfoListRequest) (int, string, interface{}, int)
	BanUser(ctx context.Context, req *req.BanUserRequest) (int, string, interface{}, int)
	EnableUser(ctx context.Context, userID string) (int, string, interface{}, int)
}

type adminServer struct {
	repo        repo.AdminRepository
	cacheCache  *cache.Cache
	redisClient *redis.Client
}

func NewAdminServer(repo repo.AdminRepository, cacheCache *cache.Cache, redisClient *redis.Client) AdminServer {
	return &adminServer{
		repo:        repo,
		cacheCache:  cacheCache,
		redisClient: redisClient,
	}
}
