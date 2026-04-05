package req

// 用户相关请求体

type RegisterRequest struct {
	Telephone string `json:"telephone" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	SmsCode   string `json:"sms_code" binding:"required"`
}

type SmsCodeRequest struct {
	Telephone string `json:"telephone" binding:"required"`
}

type LoginRequest struct {
	LoginType string `json:"login_type" binding:"required,oneof=password sms_code"`
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required_if=LoginType password"`
	Code      string `json:"code" binding:"required_if=LoginType sms_code,omitempty,len=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
