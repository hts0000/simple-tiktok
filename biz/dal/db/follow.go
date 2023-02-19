package db

import (
	"context"
	"errors"
	"log"
	"simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	FollowerID   uint   `json:"follower_id"`
	FollowerName string `json:"follower_name"`
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
func QueryFollow(ctx context.Context, uid uint) ([]*User, error) {
	follows := make([]*User, 0)
	// err := DB.WithContext(ctx).Where("follower_id = ?", uid).Find(&follows).Error
	err := DB.WithContext(ctx).Select("user.*").
		Joins("JOIN follow ON user.id = follow.user_id AND follow.follower_id = ? AND follow.deleted_at IS NULL", uid).
		Find(&follows).Error
	return follows, err
}

// uid1关注uid2
func FollowUser(ctx context.Context, uid1 uint, uid1Name string, uid2 uint, uid2Name string) error {
	follow := Follow{
		UserID:       uid2,
		Username:     uid2Name,
		FollowerID:   uid1,
		FollowerName: uid1Name,
	}
	// 先查软删除的记录
	err := DB.WithContext(ctx).Unscoped().
		Where("user_id = ? and follower_id = ?", follow.UserID, follow.FollowerID).
		Take(&follow).Error

	// 如果有，就更新deleted_at值
	if err == nil {
		// 只有删除时间不为null时才执行更新操作，避免重复触发hook函数
		if follow.DeletedAt.Valid {
			return DB.WithContext(ctx).Unscoped().
				Model(&follow).Update("deleted_at", nil).Error
		}
		return nil
	}

	// 没有软删除的记录就新加一条记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DB.WithContext(ctx).Create(&follow).Error
	}
	return err
}

// uid1取关uid2
func UnFollowUser(ctx context.Context, uid1, uid2 uint) error {
	// Follow内必须填上UserID和FollowerID，否则删除的hook中查询条件为空
	follow := Follow{
		UserID:     uid2,
		FollowerID: uid1,
	}
	err := DB.WithContext(ctx).
		Where("user_id = ? AND follower_id = ?", follow.UserID, follow.FollowerID).
		Take(&follow).Error

	// 记录存在
	if err == nil {
		if !follow.DeletedAt.Valid {
			return DB.WithContext(ctx).Delete(&follow).Error
		}
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// 查询uids列表中被uid关注
func MGetFollow(ctx context.Context, uid uint, uids []uint) ([]uint, error) {
	res := make([]uint, 0, len(uids))
	err := DB.WithContext(ctx).Model(&Follow{}).Select("user_id").Where("user_id in ? and follower_id = ?", uids, uid).Find(&res).Error
	return res, err
}

// 新增hook，当新增follow记录后，更新对应user的follow_count和follower_count
func (f *Follow) AfterCreate(tx *gorm.DB) (err error) {
	log.Println("创建hook执行")
	err = tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
	if err != nil {
		return
	}
	err = tx.Model(&User{Model: gorm.Model{ID: f.FollowerID}}).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error
	return err
}

// 更新hook，当更新follow记录后，更新对应user的follow_count和follower_count
// 因为更新只存在于更新delete_at操作，所以有更新必然是新增一条关系
func (f *Follow) AfterUpdate(tx *gorm.DB) (err error) {
	log.Println("更新hook执行")
	err = tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
	if err != nil {
		return
	}
	err = tx.Model(&User{Model: gorm.Model{ID: f.FollowerID}}).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error
	return err
}

// 删除hook，当删除follow记录后，更新对应user的follow_count和follower_count
func (f *Follow) AfterDelete(tx *gorm.DB) (err error) {
	log.Println("删除hook执行")
	err = tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error
	if err != nil {
		return
	}
	err = tx.Model(&User{Model: gorm.Model{ID: f.FollowerID}}).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error
	return err
}
