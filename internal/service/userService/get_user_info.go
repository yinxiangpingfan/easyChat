package userService

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
)

func (u *userService) GetUserInfoService(ctx context.Context) (int, string, interface{}, int) {
	log := log.FromContext(ctx)
	uuid := ctx.Value("uuid").(string)
	log.Debugf("获取用户信息, uuid: %s", uuid)
	// 从数据库中获取用户信息
	user, err := u.userRepo.GetUserInfo(ctx, uuid)
	if err != nil {
		log.Errorf("获取用户信息失败, err: %v", err)
		return errors.ErrGetUserInfoFailed.Code, errors.ErrGetUserInfoFailed.Message, nil, -1
	}
	log.Debugf("获取用户信息成功, user: %v", user)
	return 0, "", resp.UserInfoResp{
		Uuid:          user.Uuid,
		NickName:      user.NickName,
		Telephone:     user.Telephone,
		Email:         user.Email,
		Avatar:        user.Avatar,
		Signature:     user.Signature,
		Birthday:      user.Birthday,
		LastOnlineAt:  user.LastOnlineAt,
		LastOfflineAt: user.LastOfflineAt,
	}, 0
}
