package repo

import (
	"context"
	"easyChat/internal/model"
	gormscopes "easyChat/pkg/gorm_scopes"
	"errors"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetUserInfoList(ctx context.Context, page, pageSize int) ([]model.UserInfo, int64, error)
	BanUser(ctx context.Context, uuid string) error
	EnableUser(ctx context.Context, userID string) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetUserInfoList(ctx context.Context, page, pageSize int) ([]model.UserInfo, int64, error) {
	var userList []model.UserInfo
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&model.UserInfo{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Select("id", "uuid", "nickname", "telephone", "email", "avatar", "gender", "signature", "birthday", "created_at", "deleted_at", "last_online_at", "last_offline_at", "is_admin", "status").
		Scopes(gormscopes.Paginate(page, pageSize)).
		Find(&userList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return userList, total, nil
}

func (r *adminRepository) BanUser(ctx context.Context, uuid string) error {
	if err := r.db.WithContext(ctx).Model(&model.UserInfo{}).Where("uuid = ?", uuid).Update("status", 1).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminRepository) EnableUser(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Model(&model.UserInfo{}).Where("uuid = ?", userID).Update("status", 0).Error; err != nil {
		return err
	}
	return nil
}
