package db

import (
	"context"
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
