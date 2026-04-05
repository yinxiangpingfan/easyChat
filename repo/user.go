package repo

import (
	"easyChat/global"
	"easyChat/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct{}

var UserRepositoryInstance = &UserRepository{}

func (u *UserRepository) IsTelephoneRegistered(tele string) int {
	var user model.UserInfo
	if res := global.DB.Where("telephone = ?", tele).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0
		} else {
			return -1
		}
	}
	return -2
}

func (u *UserRepository) SaveUserInfo(user model.UserInfo) error {
	user.CreatedAt = time.Now()
	if res := global.DB.Create(&user); res.Error != nil {
		return res.Error
	}
	return nil
}

// GetUserByPhone 根据手机号查询用户
func (u *UserRepository) GetUserByPhone(phone string) (*model.UserInfo, error) {
	var user model.UserInfo
	if res := global.DB.Where("telephone = ?", phone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, res.Error // 查询出错
	}
	return &user, nil // 查询成功
}

// UpdateLastOnlineAt 更新用户最后登录时间
func (u *UserRepository) UpdateLastOnlineAt(uuid string) error {
	now := time.Now()
	if res := global.DB.Model(&model.UserInfo{}).Where("uuid = ?", uuid).Update("last_online_at", now); res.Error != nil {
		return res.Error
	}
	return nil
}
