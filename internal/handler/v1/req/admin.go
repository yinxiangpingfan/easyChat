package req

type UserInfoListRequest struct {
	Page     int `json:"page" binding:"required"`
	PageSize int `json:"page_size" binding:"required"`
}

type BanUserRequest struct {
	Uuid string `json:"uuid" binding:"required"`
}

// EnableUserRequest 启用用户请求
type EnableUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
