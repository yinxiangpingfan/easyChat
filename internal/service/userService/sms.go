package userService

import (
	"context"
	"crypto/rand"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/pkg/log"
	"easyChat/pkg/tools"
	"math/big"
	"time"
)

var CODE_PREFIX = "verify:code:"
var LOCK_PREFIX = "verify:lock:"
var CODE_TIMEOUT = 300
var LOCK_TIMEOUT = 60

func (u *userService) SmsCodeService(ctx context.Context, req req.SmsCodeRequest) (int, string, interface{}, int) {
	log := log.FromContext(ctx)
	isTel := tools.IsPhone(req.Telephone)
	if !isTel {
		return errors.ErrPhoneFormat.Code, errors.ErrPhoneFormat.Message, nil, -1
	}
	//判断手机号是否在验证码的限制中
	lockExist, err := u.redisClient.Exists(ctx, LOCK_PREFIX+req.Telephone).Result()
	if err != nil {
		log.Errorf("判断手机号是否在验证码的限制中失败, err: %v", err)
		return errors.ErrSmsCodeSendFailed.Code, errors.ErrSmsCodeSendFailed.Message, nil, -1
	}
	if lockExist > 0 {
		return errors.ErrSmsCodeTooFrequent.Code, errors.ErrSmsCodeTooFrequent.Message, nil, -1
	}
	//生成6位验证码
	max := big.NewInt(1000000)
	smsCode, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Errorf("生成验证码失败, err: %v", err)
		return errors.ErrSmsCodeSendFailed.Code, errors.ErrSmsCodeSendFailed.Message, nil, -1
	}
	smsCodeStr := smsCode.String()
	//发送验证码到手机号
	pipe := u.redisClient.Pipeline()
	pipe.Set(ctx, CODE_PREFIX+req.Telephone, smsCodeStr, time.Duration(CODE_TIMEOUT)*time.Second)
	pipe.Set(ctx, LOCK_PREFIX+req.Telephone, "1", time.Duration(LOCK_TIMEOUT)*time.Second)
	_, err = pipe.Exec(ctx) //把验证码存进redis
	if err != nil {
		log.Errorf("设置验证码失败, err: %v", err)
		return errors.ErrSmsCodeSendFailed.Code, errors.ErrSmsCodeSendFailed.Message, nil, -1
	}
	//发送验证码
	signName := "速通互联验证服务"       // 签名名称
	phoneNumber := req.Telephone // 手机号
	templateCode := "100001"     // 模板代码
	validDuration := "5"         // 验证码有效期（分钟）
	args := []*string{
		&signName,
		&phoneNumber,
		&templateCode,
		&validDuration,
	}
	err = tools.SendTelephoneCode(u.aliAccessConfig, args, smsCodeStr)
	if err != nil {
		//删除验证码锁
		u.redisClient.Del(ctx, LOCK_PREFIX+req.Telephone, CODE_PREFIX+req.Telephone)
		log.Errorf("发送验证码失败, err: %v", err)
		return errors.ErrSmsCodeSendFailed.Code, errors.ErrSmsCodeSendFailed.Message, nil, -1
	}
	return errors.SuccessSmsCodeSent.Code, errors.SuccessSmsCodeSent.Message, nil, 0
}
