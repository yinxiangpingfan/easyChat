package userService

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
	"encoding/json"
	"time"
)

const (
	userInfoCacheKeyPrefix = "user:info:"
	userInfoCacheTTL       = 30 * time.Minute
)

func (u *userService) GetUserInfoService(ctx context.Context) (int, string, interface{}, int) {
	logger := log.FromContext(ctx)
	uuid := ctx.Value("uuid").(string)
	cacheKey := userInfoCacheKeyPrefix + uuid

	// 1. 先从 Redis 缓存获取
	cachedData, err := u.redisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		var userInfo resp.UserInfoResp
		if err := json.Unmarshal([]byte(cachedData), &userInfo); err == nil {
			logger.Debugf("从缓存获取用户信息成功, uuid: %s", uuid)
			return 0, "", userInfo, 0
		}
	}

	// 2. 缓存未命中，从数据库获取
	user, err := u.userRepo.GetUserInfo(ctx, uuid)
	if err != nil {
		logger.Errorf("获取用户信息失败, err: %v", err)
		return errors.ErrGetUserInfoFailed.Code, errors.ErrGetUserInfoFailed.Message, nil, -1
	}

	// 3. 写入缓存
	userInfo := resp.UserInfoResp{
		Uuid:          user.Uuid,
		NickName:      user.NickName,
		Telephone:     user.Telephone,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Signature:     user.Signature,
		Birthday:      user.Birthday,
		LastOnlineAt:  user.LastOnlineAt,
		LastOfflineAt: user.LastOfflineAt,
	}
	if data, err := json.Marshal(userInfo); err == nil {
		u.redisClient.Set(ctx, cacheKey, string(data), userInfoCacheTTL)
	}

	return 0, "", userInfo, 0
}
