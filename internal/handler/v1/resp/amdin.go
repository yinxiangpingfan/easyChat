package resp

import "time"

type UserInfoListResponse struct {
	UserList []UserInfoItem `json:"user_list"`
	Total    int64          `json:"total"`
}

type UserInfoItem struct {
	Id            uint64     `json:"id"`
	Uuid          string     `json:"uuid"`
	NickName      string     `json:"nickname"`
	Telephone     string     `json:"telephone"`
	Email         string     `json:"email"`
	Avatar        string     `json:"avatar"`
	Gender        int8       `json:"gender"`
	Signature     string     `json:"signature"`
	Birthday      string     `json:"birthday"`
	CreatedAt     time.Time  `json:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	LastOnlineAt  *time.Time `json:"last_online_at"`
	LastOfflineAt *time.Time `json:"last_offline_at"`
	IsAdmin       int8       `json:"is_admin"`
	Status        int8       `json:"status"`
}
