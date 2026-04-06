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
	telephone := info[1]

	// 生成新的 Token
	accessToken, refreshToken, err := tools.GenerateAccessToken(u.jwtConfig, uuid, telephone)
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
