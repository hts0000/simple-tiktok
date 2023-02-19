package db

import (
	"context"
	"errors"
	"gorm.io/gorm"

	consts "simple-tiktok/pkg/consts"
)

type Favourite struct {
	gorm.Model
	Id      int64 `json:"id"`       //自增主键
	UserId  int64 `json:"user_id"`  //点赞用户id
	VideoId int64 `json:"video_id"` //视频id
	IsLike  int8  `json:"cancel"`   //是否点赞，0为未点赞，1为点赞

}

// TableName 修改表名映射
func (f *Favourite) TableName() string {
	return consts.FavouriteTableName
}

// IsFavourite 根据当前视频id判断是否点赞了该视频.
func IsFavourite(ctx context.Context, userId int64, videoId int64) (bool, error) {
	result := DB.WithContext(ctx).Table("likes").Unscoped().Where("user_id = ? AND video_id = ?", userId, videoId).First(&Favourite{})
	if err := result.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("can't find this data")
	} else if err != nil { //查询到数据，已经点赞过
		return true, err
	}
	return false, nil

	//未查询到数据，返回未点赞；
	/*	if result := DB.WithContext(ctx).Table("likes").Unscoped().Where("user_id = ? AND video_id = ?", userId, videoId).Take(&Favourite{}); result.RowsAffected == 0 {
			return false, errors.New("can't find this data")
		} //查询到数据
		if likeData.IsLike == consts.FavouriteAction {
			return true, nil
		} else {
			return false, nil
		}*/
}

// GetLikeUserIdList 根据videoId获取所有点赞userId
/*func GetLikeUserIdList(ctx context.Context, videoId int64) ([]int64, error) {
	var likeUserIdList []int64 //存所有该视频点赞用户id；
	//查询likes表对应视频id点赞用户，返回查询结果
	err := DB.WithContext(ctx).Model(Favourite{}).Where(map[string]interface{}{"video_id": videoId, "cancel": consts.FavouriteAction}).
		Pluck("user_id", &likeUserIdList).Error
	//查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get likeUserIdList failed")
	} else {
		//没查询到或者查询到结果，返回数量以及无报错
		return likeUserIdList, nil
	}

}*/

/*// GetLikeInfo 根据userId,videoId查询点赞信息
func GetLikeInfo(userId int64, videoId int64) (Favourite, error) {
	var favouriteInfo Favourite
	return favouriteInfo, nil
}*/

// FavouriteAction 根据userId，videoId点赞或者取消赞
func FavouriteAction(ctx context.Context, userId int64, videoId int64) error {
	//result := DB.WithContext(ctx).Unscoped().Where("user_id = ? and video_id = ?", userId, videoId).Take(&Follow{})
	result := DB.WithContext(ctx).Table("likes").Unscoped().Where("user_id = ? AND video_id = ?", userId, videoId).Take(&Favourite{})
	//点赞行为,否则取消赞
	if err := result.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return DB.WithContext(ctx).Create(&Favourite{
			UserId:  userId,
			VideoId: videoId,
		}).Error
	} else if err != nil {
		return err
	}

	// 如果有，就更新cancel为0
	return result.Update("cancel", 0).Error
}

// DisFavour 取消点赞
func DisFavour(ctx context.Context, userId int64, videoId int64) error {
	return DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Favourite{}).Error
}

// GetFavouriteList  根据userId查询like表中点赞的全部videoId
func GetFavouriteList(ctx context.Context, userId int64) ([]Video, error) {
	results := make([]Video, 0)
	err := DB.WithContext(ctx).Model(&Favourite{}).Select("user_id").Where("user_id in ? and video_id = ?", userId).Find(&results).Error
	return results, err
}

// 查找单个视频
func GetVideo(c context.Context, videoId int64) (Video, error) {
	var res Video
	result := DB.WithContext(c).Model(&Video{}).Where("id = ?", videoId).Find(&res)
	err := result.Error
	return res, err
}
