package userService

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
	"easyChat/pkg/tools"
)

func (u *userService) LoginService(ctx context.Context, req req.LoginRequest) (int, string, interface{}, int) {
	log := log.FromContext(ctx)
	// 校验手机号格式
	if !tools.IsPhone(req.Phone) {
		return errors.ErrLoginPhone.Code, errors.ErrLoginPhone.Message, nil, -1
	}

	// 查询用户
	user, err := u.userRepo.GetUserByPhone(req.Phone)
	if err != nil {
		log.Errorf("登录查询用户失败, phone: %s, err: %v", req.Phone, err)
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
		code, err := tools.VerifyCode(ctx, req.Code, "verify:code:"+req.Phone, u.redisClient)
		if err != nil {
			if code == errors.ErrSmsCodeWrong.Code {
				return errors.ErrLoginSmsCode.Code, errors.ErrLoginSmsCode.Message, nil, -1
			}
			if code == errors.ErrSmsCodeExpired.Code {
				return errors.ErrLoginSmsCode.Code, errors.ErrLoginSmsCode.Message, nil, -1
			}
			log.Errorf("登录验证码校验失败, phone: %s, err: %v", req.Phone, err)
			return errors.ErrLoginFailed.Code, errors.ErrLoginFailed.Message, nil, -1
		}
	}

	// 生成 token
	accessToken, refreshToken, err := tools.GenerateAccessToken(u.jwtConfig, user.Uuid, user.Telephone)
	if err != nil {
		log.Errorf("登录生成token失败, phone: %s, err: %v", req.Phone, err)
		return errors.ErrLoginFailed.Code, errors.ErrLoginFailed.Message, nil, -1
	}

	// 更新最后登录时间
	if err := u.userRepo.UpdateLastOnlineAt(user.Uuid); err != nil {
		log.Errorf("更新最后登录时间失败, phone: %s, err: %v", req.Phone, err)
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
