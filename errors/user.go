package errors

// ApiError API错误结构
type ApiError struct {
	Code    int
	Message string
}

// NewApiError 创建API错误
func NewApiError(code int, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
	}
}

// ==========================================验证码==========================================

var (
	ErrRequestInvalid     = NewApiError(400001, "发送验证码请求错误")
	ErrSmsCodeTooFrequent = NewApiError(400003, "发送太频繁，请稍后再试")
	ErrSmsCodeWrong       = NewApiError(400004, "验证码错误")
	ErrSmsCodeExpired     = NewApiError(400005, "验证码已过期或不存在")
	ErrSmsCodeSendFailed  = NewApiError(500001, "发送验证码失败:系统繁忙")
	SuccessSmsCodeSent    = NewApiError(200, "发送验证码成功")
)

// ==========================================注册==========================================

var (
	ErrRegisterParam  = NewApiError(400001, "参数校验失败:注册请求参数错误")
	ErrPhoneFormat    = NewApiError(400002, "手机号格式错误")
	ErrPhoneExist     = NewApiError(400003, "手机号已注册")
	ErrRegisterFailed = NewApiError(500001, "注册失败，请稍后重试")
	SuccessRegister   = NewApiError(200, "注册成功")
)
