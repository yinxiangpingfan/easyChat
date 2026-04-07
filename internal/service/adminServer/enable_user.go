package adminServer

import (
	"context"
	apiErrors "easyChat/internal/errors"
	stdErrors "errors"
	"time"

	"gorm.io/gorm"
)

func (a *adminServer) EnableUser(ctx context.Context, userID string) (int, string, interface{}, int) {
	//1.更新数据库状态
	err := a.repo.EnableUser(ctx, userID)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return apiErrors.ErrUserNotFound.Code, apiErrors.ErrUserNotFound.Message, nil, -1
		}
		return apiErrors.ErrEnableUserFailed.Code, apiErrors.ErrEnableUserFailed.Message, nil, -1
	}
	//2.更新缓存
	a.cacheCache.Set(userID, "normal", 1*time.Minute)

	//3.删除redis中的用户封禁状态
	if err := a.redisClient.Del(ctx, BannedKey+userID).Err(); err != nil {
		return apiErrors.ErrEnableUserFailed.Code, apiErrors.ErrEnableUserFailed.Message, nil, -1
	}

	return apiErrors.SuccessEnableUser.Code, apiErrors.SuccessEnableUser.Message, nil, 0
}
