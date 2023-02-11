package db

import (
	"context"
	"errors"
	"simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	UserID     uint `json:"user_id"`
	FollowerID uint `json:"follower_id"`
}

func (f *Follow) TableName() string {
	return consts.FollowTableName
}

// 查询uid的粉丝
func QueryFollower(ctx context.Context, uid uint) ([]*Follow, error) {
	followers := make([]*Follow, 0)
	err := DB.WithContext(ctx).Where("user_id = ?", uid).Find(&followers).Error
	return followers, err
}

// 查询uid的关注
func QueryFollow(ctx context.Context, uid uint) ([]*Follow, error) {
	follows := make([]*Follow, 0)
	err := DB.WithContext(ctx).Where("follower_id = ?", uid).Find(&follows).Error
	return follows, err
}

// uid1关注uid2
func FollowUser(ctx context.Context, uid1, uid2 uint) error {
	// 先查软删除的记录
	result := DB.WithContext(ctx).Unscoped().Where("user_id = ? and follower_id = ?", uid2, uid1).Take(&Follow{})
	// 没有软删除的记录就新加一条记录
	if err := result.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return DB.WithContext(ctx).Create(&Follow{
			UserID:     uid2,
			FollowerID: uid1,
		}).Error
	} else if err != nil {
		return err
	}
	// 如果有，就更新deleted_at值
	return result.Update("deleted_at", nil).Error
}

// uid1取关uid2
func UnFollowUser(ctx context.Context, uid1, uid2 uint) error {
	return DB.WithContext(ctx).Where("user_id = ? and follower_id = ?", uid2, uid1).Delete(&Follow{}).Error
}
