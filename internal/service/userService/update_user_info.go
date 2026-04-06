package userService

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
)

func (u *userService) UpdateUserInfoService(ctx context.Context, req req.UpdateUserInfoRequest) (int, string, interface{}, int) {
	log := log.FromContext(ctx)
	uuid := ctx.Value("uuid").(string)

	// 判断哪些参数需要更新
	updateFields := make(map[string]interface{})
	if req.NickName != "" {
		updateFields["nickname"] = req.NickName
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	if req.Gender != "" {
		updateFields["gender"] = req.Gender
	}
	if req.Signature != "" {
		updateFields["signature"] = req.Signature
	}
	if req.Birthday != "" {
		updateFields["birthday"] = req.Birthday
	}

	// 没有需要更新的字段
	if len(updateFields) == 0 {
		log.Warnf("更新用户信息时没有需要更新的字段, uuid: %s", uuid)
		return errors.SuccessUpdateUserInfo.Code, errors.SuccessUpdateUserInfo.Message, nil, 0
	}

	// 更新用户信息
	if err := u.userRepo.UpdateUserInfo(ctx, uuid, updateFields); err != nil {
		log.Errorf("更新用户信息失败, uuid: %s, err: %v", uuid, err)
		return errors.ErrUpdateUserInfoFailed.Code, errors.ErrUpdateUserInfoFailed.Message, nil, -1
	}

	// 清除用户信息缓存
	u.redisClient.Del(ctx, userInfoCacheKeyPrefix+uuid)
	// 返回实际更新的字段
	response := resp.UpdateUserInfoResp{
		Uuid:      uuid,
		Telephone: ctx.Value("telephone").(string),
		NickName:  req.NickName,
		Email:     req.Email,
		Gender:    req.Gender,
		Signature: req.Signature,
		Birthday:  req.Birthday,
	}
	return errors.SuccessUpdateUserInfo.Code, errors.SuccessUpdateUserInfo.Message, response, 0
}
