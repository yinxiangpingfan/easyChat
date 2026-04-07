package v1

import (
	"easyChat/internal/service/contactService"

	"github.com/gin-gonic/gin"
)

type ContactHandler interface {
	// AddContactHandler 添加联系人
	AddContactHandler() gin.HandlerFunc
	// DeleteContactHandler 删除联系人
	DeleteContactHandler() gin.HandlerFunc
	// GetContactInfoHandler 获取联系人信息
	GetContactInfoHandler() gin.HandlerFunc
}

type contactHandler struct {
	contactService contactService.ContactService
}

func NewContactHandler(contactService contactService.ContactService) ContactHandler {
	return &contactHandler{
		contactService: contactService,
	}
}

// AddContactHandler 添加联系人
func (c *contactHandler) AddContactHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: 实现添加联系人处理逻辑
	}
}

// DeleteContactHandler 删除联系人
func (c *contactHandler) DeleteContactHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: 实现删除联系人处理逻辑
	}
}

// GetContactListHandler 获取联系人列表
func (c *contactHandler) GetContactInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: 实现获取联系人列表处理逻辑
	}
}
