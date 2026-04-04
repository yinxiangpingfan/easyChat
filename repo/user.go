package repo

import (
	"easyChat/global"
	"easyChat/model"
	"errors"

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

func (u *UserRepository) SaveUserInfo(user model.UserInfo) int {
	if res := global.DB.Create(&user); res.Error != nil {
		return -1
	}
	return 0
}
