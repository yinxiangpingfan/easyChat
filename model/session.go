package model

import "time"

// 会话表
type Session struct {
	Id            uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	Uuid          string     `gorm:"uniqueIndex;type:char(41);column:uuid;comment:会话uuid"`
	SendId        string     `gorm:"type:char(41);column:send_id;index;not null;comment:创建会话人id"`
	ReceiveId     string     `gorm:"type:char(41);column:receive_id;index;not null;comment:接受会话人id"`
	ReceiveName   string     `gorm:"type:varchar(20);column:receive_name;not null;comment:名称"`
	Avatar        string     `gorm:"type:char(255);column:avatar;default:'https://ossapi.easyimpr.com/easychatava//groupDefault.jpeg';not null;comment:头像"`
	LastMessage   string     `gorm:"type:text;column:last_message;comment:最新的消息"`
	LastMessageAt *time.Time `gorm:"type:datetime;column:last_message_at;comment:最近接收时间"`
	CreatedAt     time.Time  `gorm:"type:datetime;column:created_at;index;comment:创建时间"`
	DeletedAt     *time.Time `gorm:"type:datetime;column:deleted_at;index;comment:删除时间"`
}

func (Session) TableName() string {
	return "session"
}
