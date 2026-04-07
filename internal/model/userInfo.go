package model

import "time"

// 用户信息表
type UserInfo struct {
	Id            uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	Uuid          string     `gorm:"uniqueIndex;type:char(45);column:uuid;not null;comment:用户唯一id"`
	NickName      string     `gorm:"type:varchar(20);column:nickname;not null;comment:昵称"`
	Telephone     string     `gorm:"type:char(11);column:telephone;index;not null;comment:电话"`
	Email         string     `gorm:"type:char(30);column:email;comment:邮箱"`
	Avatar        string     `gorm:"type:char(255);column:avatar;default:'https://ossapi.easyimpr.com/easychatava/default/default.png';not null;comment:头像"`
	Gender        int8       `gorm:"type:tinyint;column:gender;default:0;comment:性别(0.男,1.女)"`
	Signature     string     `gorm:"type:varchar(100);column:signature;comment:个性签名"`
	Password      string     `gorm:"type:char(32);column:password;not null;comment:哈希后的密码"`
	Birthday      string     `gorm:"type:char(8);column:birthday;comment:生日"`
	CreatedAt     time.Time  `gorm:"type:datetime;column:created_at;index;not null;comment:创建时间"`
	DeletedAt     *time.Time `gorm:"type:datetime;column:deleted_at;comment:删除时间"`
	LastOnlineAt  *time.Time `gorm:"type:datetime;column:last_online_at;comment:上次登录时间"`
	LastOfflineAt *time.Time `gorm:"type:datetime;column:last_offline_at;comment:最近离线时间"`
	IsAdmin       int8       `gorm:"type:tinyint;column:is_admin;not null;default:0;comment:是否管理员(0.否,1.是,2.是超级管理员	)"`
	Status        int8       `gorm:"type:tinyint;column:status;index;not null;default:0;comment:状态(0.正常,1.禁用)"`
	Salt          string     `gorm:"type:char(32);column:salt;not null;comment:盐值"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
