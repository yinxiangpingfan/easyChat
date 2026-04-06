package repo

import (
	"context"
	"easyChat/internal/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	IsTelephoneRegistered(ctx context.Context, tele string) int
	SaveUserInfo(ctx context.Context, user model.UserInfo) error
	GetUserByPhone(ctx context.Context, phone string) (*model.UserInfo, error)
	UpdateLastOnlineAt(ctx context.Context, uuid string) error
	GetUserInfo(ctx context.Context, uuid string) (*model.UserInfo, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// IsTelephoneRegistered 检查手机号是否已注册
func (u *userRepository) IsTelephoneRegistered(ctx context.Context, tele string) int {
	var user model.UserInfo
	if res := u.db.WithContext(ctx).Where("telephone = ?", tele).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0
		} else {
			return -1
		}
	}
	return -2
}

// SaveUserInfo 保存用户信息
func (u *userRepository) SaveUserInfo(ctx context.Context, user model.UserInfo) error {
	user.CreatedAt = time.Now()
	if res := u.db.WithContext(ctx).Create(&user); res.Error != nil {
		return res.Error
	}
	return nil
}

// GetUserByPhone 根据手机号查询用户
func (u *userRepository) GetUserByPhone(ctx context.Context, phone string) (*model.UserInfo, error) {
	var user model.UserInfo
	if res := u.db.WithContext(ctx).Where("telephone = ?", phone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, res.Error // 查询出错
	}
	return &user, nil // 查询成功
}

// GetUserInfo 根据 uuid 查询用户信息
func (u *userRepository) GetUserInfo(ctx context.Context, uuid string) (*model.UserInfo, error) {
	var user model.UserInfo
	if res := u.db.WithContext(ctx).Where("uuid = ?", uuid).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, res.Error
	}
	return &user, nil
}

// UpdateLastOnlineAt 更新用户最后登录时间
func (u *userRepository) UpdateLastOnlineAt(ctx context.Context, uuid string) error {
	now := time.Now()
	if res := u.db.WithContext(ctx).Model(&model.UserInfo{}).Where("uuid = ?", uuid).Update("last_online_at", now); res.Error != nil {
		return res.Error
	}
	return nil
}
