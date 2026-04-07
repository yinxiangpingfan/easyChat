package repo

import (
	"context"
	"easyChat/internal/model"

	"gorm.io/gorm"
)

type ContactRepository interface {
	// AddContact 添加联系人
	AddContact(ctx context.Context, contact *model.UserContact) error
	// DeleteContact 删除联系人
	DeleteContact(ctx context.Context, userId, contactId string) error
	// GetContactList 获取联系人列表
	GetContactList(ctx context.Context, userId string) ([]model.UserContact, error)
	// GetContactByUserIdAndContactId 根据用户ID和联系人ID查询联系人
	GetContactByUserIdAndContactId(ctx context.Context, userId, contactId string) (*model.UserContact, error)
}

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// AddContact 添加联系人
func (c *contactRepository) AddContact(ctx context.Context, contact *model.UserContact) error {
	// TODO: 实现添加联系人逻辑
	return nil
}

// DeleteContact 删除联系人
func (c *contactRepository) DeleteContact(ctx context.Context, userId, contactId string) error {
	// TODO: 实现删除联系人逻辑
	return nil
}

// GetContactList 获取联系人列表
func (c *contactRepository) GetContactList(ctx context.Context, userId string) ([]model.UserContact, error) {
	// TODO: 实现获取联系人列表逻辑
	return nil, nil
}

// GetContactByUserIdAndContactId 根据用户ID和联系人ID查询联系人
func (c *contactRepository) GetContactByUserIdAndContactId(ctx context.Context, userId, contactId string) (*model.UserContact, error) {
	// TODO: 实现查询联系人逻辑
	return nil, nil
}
