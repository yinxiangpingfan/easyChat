package adminServer

import (
	"context"
	apiErrors "easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/service/userService"
	"easyChat/pkg/log"
)

const (
	BannedKey = "bannedUUID:"
)

// BanUser 封禁用户
func (a *adminServer) BanUser(ctx context.Context, req *req.BanUserRequest) (int, string, interface{}, int) {
	logger := log.FromContext(ctx)

	if err := a.repo.BanUser(ctx, req.Uuid); err != nil {
		logger.Errorf("封禁用户失败, uuid: %s, err: %v", req.Uuid, err)
		return apiErrors.ErrBanUserFailed.Code, apiErrors.ErrBanUserFailed.Message, nil, -1
	}

	// 1.从redis中删除用户refreshToken
	if err := a.redisClient.Del(ctx, userService.RefreshTokenRedisKey+req.Uuid).Err(); err != nil {
		logger.Errorf("删除用户freshToken失败, uuid: %s, err: %v", req.Uuid, err)
		return apiErrors.ErrBanUserFailed.Code, apiErrors.ErrBanUserFailed.Message, nil, -1
	}

	//2.在redis中写入用户被封禁的状态
	if err := a.redisClient.Set(ctx, BannedKey+req.Uuid, "1", 0).Err(); err != nil {
		logger.Errorf("写入用户被封禁状态失败, uuid: %s, err: %v", req.Uuid, err)
		return apiErrors.ErrBanUserFailed.Code, apiErrors.ErrBanUserFailed.Message, nil, -1
	}

	//3.发布消息来让cache新增用户封禁
	if err := a.redisClient.Publish(ctx, "channel:user_ban", req.Uuid).Err(); err != nil {
		logger.Errorf("发布封禁用户消息失败, uuid: %s, err: %v", req.Uuid, err)
		return apiErrors.ErrBanUserFailed.Code, apiErrors.ErrBanUserFailed.Message, nil, -1
	}

	logger.Infof("封禁用户成功, uuid: %s", req.Uuid)
	return apiErrors.SuccessBanUser.Code, apiErrors.SuccessBanUser.Message, nil, 0
}
