package errors

// ======================================获取用户列表======================================
var (
	ErrGetUserListFailed = ApiError{Code: 500001, Message: "系统繁忙"}
	SuccessGetUserList   = ApiError{Code: 200, Message: "获取用户列表成功"}
)

// ======================================封禁用户======================================
var (
	ErrBanUserFailed = ApiError{Code: 500001, Message: "系统繁忙"}
	SuccessBanUser   = ApiError{Code: 200, Message: "封禁用户成功"}
)

// ======================================启用用户======================================
var (
	ErrEnableUserFailed   = ApiError{Code: 500002, Message: "启用用户失败"}
	ErrUserNotFound       = ApiError{Code: 500003, Message: "用户不存在"}
	ErrUserAlreadyEnabled = ApiError{Code: 500004, Message: "用户已是启用状态"}
	SuccessEnableUser     = ApiError{Code: 200, Message: "启用用户成功"}
)

// ======================================删除用户======================================
var (
	ErrDeleteUserNotFound = ApiError{Code: 500101, Message: "用户不存在"}
	ErrDeleteUserFailed   = ApiError{Code: 500102, Message: "删除用户失败"}
	SuccessDeleteUser     = ApiError{Code: 200, Message: "删除用户成功"}
)
