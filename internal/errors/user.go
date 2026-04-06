package errors

// ApiError API错误结构
type ApiError struct {
	Code    int
	Message string
}

// ==========================================验证码==========================================

var (
	ErrRequestInvalid     = ApiError{Code: 400001, Message: "发送验证码请求错误"}
	ErrSmsCodeTooFrequent = ApiError{Code: 400003, Message: "发送太频繁，请稍后再试"}
	ErrSmsCodeWrong       = ApiError{Code: 400004, Message: "验证码错误"}
	ErrSmsCodeExpired     = ApiError{Code: 400005, Message: "验证码已过期或不存在"}
	ErrSmsCodeSendFailed  = ApiError{Code: 500001, Message: "发送验证码失败:系统繁忙"}
	SuccessSmsCodeSent    = ApiError{Code: 200, Message: "发送验证码成功"}
)

// ==========================================注册==========================================

var (
	ErrRegisterParam  = ApiError{Code: 400001, Message: "参数校验失败:注册请求参数错误"}
	ErrPhoneFormat    = ApiError{Code: 400002, Message: "手机号格式错误"}
	ErrPhoneExist     = ApiError{Code: 400003, Message: "手机号已注册"}
	ErrRegisterFailed = ApiError{Code: 500001, Message: "注册失败，请稍后重试"}
	SuccessRegister   = ApiError{Code: 200, Message: "注册成功"}
)

// ==========================================鉴权==========================================

var (
	ErrAuthFormat  = ApiError{Code: 401000, Message: "AccessToken格式错误"}
	ErrAuthExpired = ApiError{Code: 401001, Message: "AccessToken已过期或无效"}
	ErrAuthFailed  = ApiError{Code: 500001, Message: "鉴权失败，请稍后重试"}
)

// ======================================登录======================================
var (
	ErrLoginParam        = ApiError{Code: 400001, Message: "参数校验失败:登录请求参数错误"}
	ErrLoginPhone        = ApiError{Code: 400002, Message: "手机号格式错误"}
	ErrLoginUserNotExist = ApiError{Code: 400003, Message: "用户不存在"}
	ErrLoginPassword     = ApiError{Code: 400004, Message: "密码错误"}
	ErrLoginSmsCode      = ApiError{Code: 400005, Message: "验证码错误或已过期"}
	ErrLoginFailed       = ApiError{Code: 500001, Message: "登录失败，请稍后重试"}
	SuccessLogin         = ApiError{Code: 200, Message: "登录成功"}
)

// ======================================刷新Token======================================
var (
	ErrRefreshTokenInvalid       = ApiError{Code: 401002, Message: "RefreshToken无效或已过期"}
	SuccessRefreshToken          = ApiError{Code: 200, Message: "刷新Token成功"}
	ErrRefreshTokenRefreshFailed = ApiError{Code: 500001, Message: "刷新Token时失败，请稍后重试"}
)
