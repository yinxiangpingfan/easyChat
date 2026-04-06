package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 校验验证码是否正确
func VerifyCode(c context.Context, verifyCode, key string, redisClient *redis.Client) (int, error) {
	//从redis中获取验证码
	smsCode, err := redisClient.Get(c, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 400005, fmt.Errorf("验证码已过期或不存在")
		}
		return 500001, fmt.Errorf("系统繁忙")
	}
	if smsCode != verifyCode {
		return 400004, fmt.Errorf("验证码错误")
	}
	redisClient.Del(c, key)
	return 0, nil
}
