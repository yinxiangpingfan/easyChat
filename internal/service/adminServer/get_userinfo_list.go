package adminServer

import (
	"context"
	"easyChat/internal/errors"
	"easyChat/internal/handler/v1/req"
	"easyChat/internal/handler/v1/resp"
	"easyChat/pkg/log"
)

// GetUserInfoList 获取用户列表
// admin接口，不做redis缓存
func (a *adminServer) GetUserInfoList(ctx context.Context, req *req.UserInfoListRequest) (int, string, interface{}, int) {
	logger := log.FromContext(ctx)

	userList, total, err := a.repo.GetUserInfoList(ctx, req.Page, req.PageSize)
	if err != nil {
		logger.Errorf("获取用户列表失败, err: %v", err)
		return errors.ErrGetUserListFailed.Code, errors.ErrGetUserListFailed.Message, nil, -1
	}
	if userList == nil {
		return errors.SuccessGetUserList.Code, errors.SuccessGetUserList.Message, resp.UserInfoListResponse{
			UserList: []resp.UserInfoItem{},
			Total:    0,
		}, 0
	}

	// 转换为响应体
	items := make([]resp.UserInfoItem, 0, len(userList))
	for _, user := range userList {
		items = append(items, resp.UserInfoItem{
			Id:            user.Id,
			Uuid:          user.Uuid,
			NickName:      user.NickName,
			Telephone:     user.Telephone,
			Email:         user.Email,
			Avatar:        user.Avatar,
			Gender:        user.Gender,
			Signature:     user.Signature,
			Birthday:      user.Birthday,
			CreatedAt:     user.CreatedAt,
			DeletedAt:     user.DeletedAt,
			LastOnlineAt:  user.LastOnlineAt,
			LastOfflineAt: user.LastOfflineAt,
			IsAdmin:       user.IsAdmin,
			Status:        user.Status,
		})
	}

	return errors.SuccessGetUserList.Code, errors.SuccessGetUserList.Message, resp.UserInfoListResponse{
		UserList: items,
		Total:    total,
	}, 0
}
