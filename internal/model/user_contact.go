package model

import "time"

// 用户联系人表
type UserContact struct {
	Id          uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	UserId      string     `gorm:"type:char(20);column:user_id;index;not null;comment:用户唯一id"`
	ContactId   string     `gorm:"type:char(20);column:contact_id;index;not null;comment:对应联系id"`
	ContactType int8       `gorm:"type:tinyint;column:contact_type;not null;comment:联系类型(0.用户,1.群聊)"`
	Status      int8       `gorm:"type:tinyint;column:status;not null;comment:联系状态"`
	CreatedAt   time.Time  `gorm:"type:datetime;column:created_at;not null;comment:创建时间"`
	UpdatedAt   time.Time  `gorm:"type:datetime;column:update_at;not null;comment:更新时间"`
	DeletedAt   *time.Time `gorm:"type:datetime;column:deleted_at;index;comment:删除时间"`
}

func (UserContact) TableName() string {
	return "user_contact"
}
