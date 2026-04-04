package model

import "time"

// 消息表
type Message struct {
	Id         uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	Uuid       string     `gorm:"uniqueIndex;type:char(41);column:uuid;not null;comment:消息uuid"`
	SessionId  string     `gorm:"type:char(41);column:session_id;index;not null;comment:会话uuid"`
	Type       int8       `gorm:"type:tinyint;column:type;not null;comment:消息类型(0.文本,1.语音,2.文件,3.通话)"`
	Content    string     `gorm:"type:text;column:content;comment:消息内容"`
	Url        string     `gorm:"type:char(255);column:url;comment:消息url"`
	SendId     string     `gorm:"type:char(41);column:send_id;index;not null;comment:发送者uuid"`
	SendName   string     `gorm:"type:varchar(20);column:send_name;not null;comment:发送者昵称"`
	SendAvatar string     `gorm:"type:varchar(255);column:send_avatar;not null;comment:发送者头像"`
	ReceiveId  string     `gorm:"type:char(41);column:receive_id;index;not null;comment:接受者uuid"`
	FileType   string     `gorm:"type:char(10);column:file_type;comment:文件类型"`
	FileName   string     `gorm:"type:varchar(50);column:file_name;comment:文件名"`
	FileSize   string     `gorm:"type:char(20);column:file_size;comment:文件大小"`
	Status     int8       `gorm:"type:tinyint;column:status;not null;comment:状态(0.未发送,1.已发送)"`
	CreatedAt  time.Time  `gorm:"type:datetime;column:created_at;not null;comment:创建时间"`
	SendAt     *time.Time `gorm:"type:datetime;column:send_at;comment:发送时间"`
	AvData     string     `gorm:"type:text;column:av_data;comment:通话传递数据"`
}

func (Message) TableName() string {
	return "message"
}
