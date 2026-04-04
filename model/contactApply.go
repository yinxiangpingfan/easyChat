package model

import "time"

// 联系人申请表
type ContactApply struct {
	Id          uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	Uuid        string     `gorm:"uniqueIndex;type:char(41);column:uuid;comment:申请id"`
	UserId      string     `gorm:"type:char(20);column:user_id;index;not null;comment:申请人id"`
	ContactId   string     `gorm:"type:char(20);column:contact_id;index;not null;comment:被申请id"`
	ContactType int8       `gorm:"type:tinyint;column:contact_type;not null;comment:被申请类型(0.用户,1.群聊)"`
	Status      int8       `gorm:"type:tinyint;column:status;not null;comment:申请状态"`
	Message     string     `gorm:"type:varchar(100);column:message;comment:申请信息"`
	LastApplyAt time.Time  `gorm:"type:datetime;column:last_apply_at;not null;comment:最后申请时间"`
	DeletedAt   *time.Time `gorm:"type:datetime;column:deleted_at;index;comment:删除时间"`
}

func (ContactApply) TableName() string {
	return "contact_apply"
}
