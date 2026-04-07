package userService

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
	"easyChat/pkg/tools"
)

func (u *userService) RefreshTokenService(ctx context.Context, req req.RefreshTokenRequest) (int, string, interface{}, int) {
	log := log.FromContext(ctx)
	// 验证 RefreshToken
	info, code, err := tools.VerifyRefreshToken(u.jwtConfig, req.RefreshToken)
	if err != nil {
		if code == 2 {
			log.Errorf("RefreshToken验证时发生系统性失败: %v", err)
			return errors.ErrRefreshTokenRefreshFailed.Code, errors.ErrRefreshTokenRefreshFailed.Message, nil, -1
		}
		return errors.ErrRefreshTokenInvalid.Code, errors.ErrRefreshTokenInvalid.Message, nil, -1
	}

	uuid := info[0]
	//检验redis中是否有对应的refreshToken
	exists, err := u.redisClient.Exists(ctx, RefreshTokenRedisKey+uuid).Result()
	if err != nil {
		log.Errorf("RefreshToken验证时发生系统性失败: %v", err)
		return errors.ErrRefreshTokenRefreshFailed.Code, errors.ErrRefreshTokenRefreshFailed.Message, nil, -1
	}
	if exists == 0 {
		return errors.ErrRefreshTokenInvalid.Code, errors.ErrRefreshTokenInvalid.Message, nil, -1
	}
	telephone := info[1]
	isAdmin := 0
	switch info[2] {
	case "admin":
		isAdmin = 1
	case "superAdmin":
		isAdmin = 2
	default:
		isAdmin = 0
	}

	// 生成新的 Token
	accessToken, refreshToken, err := tools.GenerateAccessToken(u.jwtConfig, uuid, telephone, int8(isAdmin))
	if err != nil {
		log.Errorf("刷新Token失败: %v", err)
		return errors.ErrRefreshTokenRefreshFailed.Code, errors.ErrRefreshTokenRefreshFailed.Message, nil, -1
	}

	// 返回新的 Token
	response := resp.RefreshTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return errors.SuccessRefreshToken.Code, errors.SuccessRefreshToken.Message, response, 0
}
