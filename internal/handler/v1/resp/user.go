package resp

import (
	"time"
)

type RegisterResp struct {
	Uuid      string `json:"uuid"`
	NickName  string `json:"nick_name"`
	Telephone string `json:"telephone"`
}

type LoginResp struct {
	Uuid         string `json:"uuid"`
	NickName     string `json:"nick_name"`
	Telephone    string `json:"telephone"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfoResp struct {
	Uuid          string
	NickName      string
	Telephone     string
	Email         string
	Avatar        string
	Signature     string
	Birthday      string
	LastOnlineAt  *time.Time
	LastOfflineAt *time.Time
}

type UpdateUserInfoResp struct {
	Uuid      string `json:"uuid"`
	Telephone string `json:"telephone,omitempty"`
	NickName  string `json:"nickname,omitempty"`
	Email     string `json:"email,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Signature string `json:"signature,omitempty"`
	Birthday  string `json:"birthday,omitempty"`
}
