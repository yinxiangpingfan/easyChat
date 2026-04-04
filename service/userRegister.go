package service

import (
	"easyChat/global"
	"easyChat/handler/v1/req"
	"easyChat/model"
	"easyChat/repo"
	"easyChat/tools"
	"time"

	"github.com/google/uuid"
)

func (u *UserService) RegisterService(req req.RegisterRequest) (int, string, interface{}, int) {
	//校验参数
	if tools.IsPhone(req.Telephone) {
		global.Logger.Debugf("手机号格式错误, req: %v", req.Telephone)
		return 400002, "手机号格式错误", nil, -1
	}
	//校验验证码

	//判断手机号是否注册过
	res := repo.UserRepositoryInstance.IsTelephoneRegistered(req.Telephone)
	switch res {
	case -1:
		global.Logger.Errorf("注册时查询数据库失败, req: %v", req.Telephone)
		return 500001, "注册失败，请稍后重试", nil, -1
	case -2:
		global.Logger.Infof("注册时手机号已注册, req: %v", req.Telephone)
		return 400003, "手机号已注册", nil, -1
	}
	//加密密码
	salt := tools.GenerateSalt()
	if salt == "" {
		global.Logger.Errorf("注册时生成盐值失败, req: %v", req.Telephone)
		return 500001, "注册失败，请稍后重试", nil, -1
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
		return 500001, "注册失败，请稍后重试", nil, -1
	}
	global.Logger.Infof("注册时保存用户信息到数据库成功, req: %v", req.Telephone)
	//返回注册成功
	response := map[string]string{
		"telephone": req.Telephone,
		"uuid":      uuid,
		"nickname":  req.Nickname,
	}
	return 200, "注册成功", response, 0
}
