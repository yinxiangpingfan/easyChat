package repo

import (
	"easyChat/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	IsTelephoneRegistered(tele string) int
	SaveUserInfo(user model.UserInfo) error
	GetUserByPhone(phone string) (*model.UserInfo, error)
	UpdateLastOnlineAt(uuid string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// IsTelephoneRegistered 检查手机号是否已注册
func (u *userRepository) IsTelephoneRegistered(tele string) int {
	var user model.UserInfo
	if res := u.db.Where("telephone = ?", tele).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0
		} else {
			return -1
		}
	}
	return -2
}

// SaveUserInfo 保存用户信息
func (u *userRepository) SaveUserInfo(user model.UserInfo) error {
	user.CreatedAt = time.Now()
	if res := u.db.Create(&user); res.Error != nil {
		return res.Error
	}
	return nil
}

// GetUserByPhone 根据手机号查询用户
func (u *userRepository) GetUserByPhone(phone string) (*model.UserInfo, error) {
	var user model.UserInfo
	if res := u.db.Where("telephone = ?", phone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, res.Error // 查询出错
	}
	return &user, nil // 查询成功
}

// UpdateLastOnlineAt 更新用户最后登录时间
func (u *userRepository) UpdateLastOnlineAt(uuid string) error {
	now := time.Now()
	if res := u.db.Model(&model.UserInfo{}).Where("uuid = ?", uuid).Update("last_online_at", now); res.Error != nil {
		return res.Error
	}
	return nil
}
