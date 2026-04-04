package model

import "time"

// 群组表
type GroupInfo struct {
	Id        uint64     `gorm:"primaryKey;column:id;comment:自增id"`
	Uuid      string     `gorm:"uniqueIndex;type:char(41);column:uuid;not null;comment:群组唯一id"`
	Name      string     `gorm:"type:varchar(20);column:name;not null;comment:群名称"`
	Notice    string     `gorm:"type:varchar(500);column:notice;comment:群公告"`
	Members   string     `gorm:"type:json;column:members;comment:群组成员"`
	MemberCnt int        `gorm:"type:int;column:member_cnt;default:1;comment:群人数"`
	OwnerId   string     `gorm:"type:char(20);column:owner_id;not null;comment:群主uuid"`
	AddMode   int8       `gorm:"type:tinyint;column:add_mode;default:0;comment:加群方式(0.直接,1.审核)"`
	Avatar    string     `gorm:"type:char(255);column:avatar;default:'https://ossapi.easyimpr.com/easychatava/default/groupDefault.jpeg';not null;comment:头像"`
	Status    int8       `gorm:"type:tinyint;column:status;default:0;comment:状态(0.正常,1.禁用,2.解散)"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;index;not null;comment:创建时间"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;not null;comment:更新时间"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at;index;comment:删除时间"`
}

func (GroupInfo) TableName() string {
	return "group_info"
}
