package v1

import (
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/service/adminServer"
	"easyChat/pkg/log"

	"github.com/gin-gonic/gin"
)

type AdminHandler interface {
	GetUserInfoList() gin.HandlerFunc
	BanUser() gin.HandlerFunc
	EnableUser() gin.HandlerFunc
}

type adminHandler struct {
	service adminServer.AdminServer
}

func NewAdminHandler(service adminServer.AdminServer) AdminHandler {
	return &adminHandler{service: service}
}

func (a *adminHandler) BanUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &req.BanUserRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Warnf("绑定请求参数发生了错误 err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, msg, data, ret := a.service.BanUser(c, req)
		JsonBack(c, code, msg, data, ret)
	}
}

func (a *adminHandler) GetUserInfoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &req.UserInfoListRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Warnf("绑定请求参数发生了错误 err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, msg, data, ret := a.service.GetUserInfoList(c, req)
		JsonBack(c, code, msg, data, ret)
	}
}

func (a *adminHandler) EnableUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &req.EnableUserRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			reqLog := log.FromContext(c.Request.Context())
			reqLog.Warnf("绑定请求参数发生了错误 err: %v", err)
			JsonBack(c, errors.ErrRequestInvalid.Code, errors.ErrRequestInvalid.Message, nil, -1)
			return
		}
		code, msg, data, ret := a.service.EnableUser(c, req.UserID)
		JsonBack(c, code, msg, data, ret)
	}
}

// DeleteUser 删除用户
