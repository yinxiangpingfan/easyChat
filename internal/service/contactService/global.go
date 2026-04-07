package contactService

import (
	"context"
	"easyChat/internal/repo"

	"github.com/redis/go-redis/v9"
)

type ContactService interface {
	// AddContactService 添加联系人
	AddContactService(ctx context.Context, userId, contactId string, contactType int8) (int, string, interface{}, int)
	// DeleteContactService 删除联系人
	DeleteContactService(ctx context.Context, userId, contactId string) (int, string, interface{}, int)
	// GetContactListService 获取联系人列表
	GetContactListService(ctx context.Context, userId string) (int, string, interface{}, int)
}

type contactService struct {
	contactRepo repo.ContactRepository
	redisClient *redis.Client
}

func NewContactService(contactRepo repo.ContactRepository, redisClient *redis.Client) ContactService {
	return &contactService{
		contactRepo: contactRepo,
		redisClient: redisClient,
	}
}

// AddContactService 添加联系人
func (c *contactService) AddContactService(ctx context.Context, userId, contactId string, contactType int8) (int, string, interface{}, int) {
	// TODO: 实现添加联系人业务逻辑
	return 0, "", nil, 0
}

// DeleteContactService 删除联系人
func (c *contactService) DeleteContactService(ctx context.Context, userId, contactId string) (int, string, interface{}, int) {
	// TODO: 实现删除联系人业务逻辑
	return 0, "", nil, 0
}

// GetContactListService 获取联系人列表
func (c *contactService) GetContactListService(ctx context.Context, userId string) (int, string, interface{}, int) {
	// TODO: 实现获取联系人列表业务逻辑
	return 0, "", nil, 0
}
