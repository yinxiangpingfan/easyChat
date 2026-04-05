package userService

import (
	"easyChat/errors"
	"easyChat/global"
	"easyChat/handler/v1/req"
	"easyChat/handler/v1/resp"
	"easyChat/repo"
	"easyChat/tools"

	"github.com/gin-gonic/gin"
)

func (u *UserService) LoginService(c *gin.Context, req req.LoginRequest) (int, string, interface{}, int) {
	// 校验手机号格式
	if !tools.IsPhone(req.Phone) {
		return errors.ErrLoginPhone.Code, errors.ErrLoginPhone.Message, nil, -1
	}

	// 查询用户
	user, err := repo.UserRepositoryInstance.GetUserByPhone(req.Phone)
	if err != nil {
		global.Logger.Errorf("登录查询用户失败, phone: %s, err: %v", req.Phone, err)
		return errors.ErrLoginFailed.Code, errors.ErrLoginFailed.Message, nil, -1
	}
	if user == nil {
		return errors.ErrLoginUserNotExist.Code, errors.ErrLoginUserNotExist.Message, nil, -1
	}

	if req.LoginType == "password" {
		// 密码登录
		hashedPassword := tools.PasswordHash(req.Password, user.Salt)
		if hashedPassword != user.Password {
			return errors.ErrLoginPassword.Code, errors.ErrLoginPassword.Message, nil, -1
		}
	} else {
		// 验证码登录
		code, err := tools.VerifyCode(c, req.Code, "verify:code:"+req.Phone)
		if err != nil {
			if code == errors.ErrSmsCodeWrong.Code {
				return errors.ErrLoginSmsCode.Code, errors.ErrLoginSmsCode.Message, nil, -1
			}
			if code == errors.ErrSmsCodeExpired.Code {
				return errors.ErrLoginSmsCode.Code, errors.ErrLoginSmsCode.Message, nil, -1
			}
			global.Logger.Errorf("登录验证码校验失败, phone: %s, err: %v", req.Phone, err)
			return errors.ErrLoginFailed.Code, errors.ErrLoginFailed.Message, nil, -1
		}
	}

	// 生成 token
	accessToken, refreshToken, err := tools.GenerateAccessToken(user.Uuid, user.Telephone)
	if err != nil {
		global.Logger.Errorf("登录生成token失败, phone: %s, err: %v", req.Phone, err)
		return errors.ErrLoginFailed.Code, errors.ErrLoginFailed.Message, nil, -1
	}

	// 返回登录成功
	response := resp.LoginResp{
		Uuid:         user.Uuid,
		NickName:     user.NickName,
		Telephone:    user.Telephone,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	return errors.SuccessLogin.Code, errors.SuccessLogin.Message, response, 0
}
