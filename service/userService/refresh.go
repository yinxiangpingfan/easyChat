package userService

import (
	"easyChat/errors"
	"easyChat/global"
	"easyChat/handler/v1/req"
	"easyChat/handler/v1/resp"
	"easyChat/tools"

	"github.com/gin-gonic/gin"
)

func (u *userService) RefreshTokenService(c *gin.Context, req req.RefreshTokenRequest) (int, string, interface{}, int) {
	// 验证 RefreshToken
	info, code, err := tools.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		if code == 2 {
			global.Logger.Errorf("RefreshToken验证时发生系统性失败: %v", err)
			return errors.ErrRefreshTokenRefreshFailed.Code, errors.ErrRefreshTokenRefreshFailed.Message, nil, -1
		}
		return errors.ErrRefreshTokenInvalid.Code, errors.ErrRefreshTokenInvalid.Message, nil, -1
	}

	uuid := info[0]
	telephone := info[1]

	// 生成新的 Token
	accessToken, refreshToken, err := tools.GenerateAccessToken(uuid, telephone)
	if err != nil {
		global.Logger.Errorf("刷新Token失败: %v", err)
		return errors.ErrRefreshTokenRefreshFailed.Code, errors.ErrRefreshTokenRefreshFailed.Message, nil, -1
	}

	// 返回新的 Token
	response := resp.RefreshTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return errors.SuccessRefreshToken.Code, errors.SuccessRefreshToken.Message, response, 0
}
