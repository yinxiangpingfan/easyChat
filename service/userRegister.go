package service

import (
	"context"
	"easyChat/errors"
	"easyChat/global"
	"easyChat/handler/v1/req"
	"easyChat/handler/v1/resp"
	"easyChat/model"
	"easyChat/repo"
	"easyChat/tools"
	"time"

	"github.com/google/uuid"
)

func (u *UserService) RegisterService(ctx context.Context, req req.RegisterRequest) (int, string, interface{}, int) {
	//校验参数
	if !tools.IsPhone(req.Telephone) {
		global.Logger.Debugf("手机号格式错误, req: %v", req.Telephone)
		return errors.ErrPhoneFormat.Code, errors.ErrPhoneFormat.Message, nil, -1
	}
	//校验验证码
	code, err := tools.VerifyCode(ctx, req.SmsCode, CODE_PREFIX+req.Telephone)
	if err != nil {
		switch code {
		case errors.ErrSmsCodeWrong.Code:
			return errors.ErrSmsCodeWrong.Code, errors.ErrSmsCodeWrong.Message, nil, -1
		case errors.ErrSmsCodeExpired.Code:
			return errors.ErrSmsCodeExpired.Code, errors.ErrSmsCodeExpired.Message, nil, -1
		case errors.ErrRegisterFailed.Code:
			global.Logger.Errorf("注册时校验验证码失败, 发生错误, req: %v", req.Telephone)
			return errors.ErrRegisterFailed.Code, errors.ErrRegisterFailed.Message, nil, -1
		}
	}

	//判断手机号是否注册过
	res := repo.UserRepositoryInstance.IsTelephoneRegistered(req.Telephone)
	switch res {
	case -1:
		global.Logger.Errorf("注册时查询数据库失败, req: %v", req.Telephone)
		return errors.ErrRegisterFailed.Code, errors.ErrRegisterFailed.Message, nil, -1
	case -2:
		global.Logger.Infof("注册时手机号已注册, req: %v", req.Telephone)
		return errors.ErrPhoneExist.Code, errors.ErrPhoneExist.Message, nil, -1
	}
	//加密密码
	salt := tools.GenerateSalt()
	if salt == "" {
		global.Logger.Errorf("注册时生成盐值失败, req: %v", req.Telephone)
		return errors.ErrRegisterFailed.Code, errors.ErrRegisterFailed.Message, nil, -1
	}
	req.Password = tools.PasswordHash(req.Password, salt)
	//保存用户信息到数据库
	uuid := "U" + time.Now().Format("20060102") + uuid.New().String() //U+年月日+uuid
	res = repo.UserRepositoryInstance.SaveUserInfo(model.UserInfo{
		Uuid:      uuid,
		Telephone: req.Telephone,
		Password:  req.Password,
		NickName:  req.Nickname,
		Salt:      salt,
	})
	switch res {
	case -1:
		global.Logger.Errorf("注册时保存用户信息到数据库失败, req: %v", req.Telephone)
		return errors.ErrRegisterFailed.Code, errors.ErrRegisterFailed.Message, nil, -1
	}
	global.Logger.Infof("注册时保存用户信息到数据库成功, req: %v", req.Telephone)
	//返回注册成功
	response := resp.RegisterResp{
		Uuid:      uuid,
		NickName:  req.Nickname,
		Telephone: req.Telephone,
	}
	return errors.SuccessRegister.Code, errors.SuccessRegister.Message, response, 0
}
