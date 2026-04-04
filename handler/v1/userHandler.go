package v1

import (
	"easyChat/global"
	"easyChat/handler/v1/req"
	"easyChat/service"

	"github.com/gin-gonic/gin"
)

// 发送验证码
func SmsCodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.SmsCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			global.Logger.Warnf("发送验证码时参数绑定失败, err: %v", err)
			JsonBack(c, 400001, "发送验证码请求错误", nil, -1)
			return
		}
		code, message, data, ret := service.UserServiceInstance.SmsCodeService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}

// 注册
func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			global.Logger.Warnf("注册时参数绑定失败, err: %v", err)
			JsonBack(c, 400001, "参数校验失败:注册请求参数错误", nil, -1)
			return
		}
		code, message, data, ret := service.UserServiceInstance.RegisterService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}
