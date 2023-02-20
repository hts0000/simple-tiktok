package db

import (
	"context"
	"errors"
	"log"
	consts "simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserID  uint  `json:"user_id"`  //点赞用户id
	VideoID uint  `json:"video_id"` //视频id
	Cancel  uint8 `json:"cancel"`   // 默认不点赞为0，点赞为1
}

// TableName 修改表名映射
func (f *Favorite) TableName() string {
	return consts.FavouriteTableName
}

// IsFavourite 根据当前视频id判断是否点赞了该视频.
func IsLike(ctx context.Context, videoId int64, userId int64) (bool, error) {
	like := Favorite{}
	//未查询到数据，返回未点赞
	err := DB.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Take(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// uid点赞vid
func LikeVideo(ctx context.Context, uid, vid uint) error {
	fav := Favorite{
		UserID:  uid,
		VideoID: vid,
		Cancel:  1,
	}
	// 先查软删除的记录
	err := DB.WithContext(ctx).Unscoped().
		Where("user_id = ? AND video_id = ?", fav.UserID, fav.VideoID).
		Take(&fav).Error

	// 如果有，就更新deleted_at值
	if err == nil {
		// 只有删除时间不为null时才执行更新操作，避免重复触发hook函数
		if fav.DeletedAt.Valid {
			return DB.WithContext(ctx).Unscoped().
				Model(&fav).Update("deleted_at", nil).Update("cancel", fav.Cancel).Error
		}
		return nil
	}

	// 没有软删除的记录就新加一条记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return DB.WithContext(ctx).Create(&fav).Error
	}
	return err
}

// uid取消点赞vid
func UnLikeVideo(ctx context.Context, uid, vid uint) error {
	fav := Favorite{
		UserID:  uid,
		VideoID: vid,
		Cancel:  0,
	}
	err := DB.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", fav.UserID, fav.VideoID).
		Take(&fav).Error

	// 记录存在
	if err == nil {
		if !fav.DeletedAt.Valid {
			return DB.WithContext(ctx).Delete(&fav).Error
		}
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

// 获取like的视频列表
func GetLikes(ctx context.Context, uid uint) ([]*Video, error) {
	videos := make([]*Video, 0)
	err := DB.WithContext(ctx).Select("videos.*").
		Joins("JOIN likes ON likes.video_id = videos.id AND likes.user_id = ? AND likes.deleted_at IS NULL", uid).
		Find(&videos).Error
	return videos, err
}

// 新增hook，当新增like记录后，更新对应user的favorite_count
func (f *Favorite) AfterCreate(tx *gorm.DB) (err error) {
	log.Println("创建hook执行")
	return tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

// 更新hook，当更新like记录后，更新对应user的favorite_count
// 因为更新只存在于更新delete_at操作，所以有更新必然是新增一条关系
func (f *Favorite) AfterUpdate(tx *gorm.DB) (err error) {
	log.Println("更新hook执行")
	return tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

// 删除hook，当删除like记录后，更新对应user的favorite_count
func (f *Favorite) AfterDelete(tx *gorm.DB) (err error) {
	log.Println("删除hook执行")
	return tx.Model(&User{Model: gorm.Model{ID: f.UserID}}).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}
