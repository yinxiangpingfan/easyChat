package v1

import (
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/service/userService"
	"easyChat/pkg/log"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	SmsCodeHandler() gin.HandlerFunc
	RegisterHandler() gin.HandlerFunc
	LoginHandler() gin.HandlerFunc
	RefreshTokenHandler() gin.HandlerFunc
	GetUserInfoHandler() gin.HandlerFunc
	UpdateUserInfoHandler() gin.HandlerFunc
}
type userHandler struct {
	userService userService.UserService
}

func NewUserHandler(userService userService.UserService) UserHandler {
	return &userHandler{userService: userService}
}

// 发送验证码
func (u *userHandler) SmsCodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.SmsCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.FromContext(c.Request.Context()).Warnf("发送验证码时参数绑定失败, err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, message, data, ret := u.userService.SmsCodeService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}

// 注册
func (u *userHandler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.FromContext(c.Request.Context()).Warnf("注册时参数绑定失败, err: %v", err)
			JsonBack(c, errors.ErrRegisterParam.Code, errors.ErrRegisterParam.Message, nil, -1)
			return
		}
		code, message, data, ret := u.userService.RegisterService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}

// 登录
func (u *userHandler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.FromContext(c.Request.Context()).Warnf("登录时参数绑定失败, err: %v", err)
			JsonBack(c, errors.ErrLoginParam.Code, errors.ErrLoginParam.Message, nil, -1)
			return
		}
		code, message, data, ret := u.userService.LoginService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}

// 刷新 Token
func (u *userHandler) RefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.FromContext(c.Request.Context()).Warnf("刷新 Token时参数绑定失败, err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, message, data, ret := u.userService.RefreshTokenService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}

// 获取用户信息
func (u *userHandler) GetUserInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		code, message, data, ret := u.userService.GetUserInfoService(c)
		JsonBack(c, code, message, data, ret)
	}
}

// 更新用户信息

func (u *userHandler) UpdateUserInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req req.UpdateUserInfoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.FromContext(c.Request.Context()).Warnf("更新用户信息时参数绑定失败, err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, message, data, ret := u.userService.UpdateUserInfoService(c, req)
		JsonBack(c, code, message, data, ret)
	}
}
