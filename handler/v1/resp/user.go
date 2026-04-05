package resp

type RegisterResp struct {
	Uuid      string `json:"uuid"`
	NickName  string `json:"nick_name"`
	Telephone string `json:"telephone"`
}
