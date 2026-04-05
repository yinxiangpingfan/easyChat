package resp

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
