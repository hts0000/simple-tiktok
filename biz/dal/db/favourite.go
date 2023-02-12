package db

import (
	"gorm.io/gorm"
	"simple-tiktok/pkg/consts"
)

type Favourite struct {
	gorm.Model
	Id      int64 `json:"id"`       //自增主键
	UserId  int64 `json:"user_id"`  //点赞用户id
	VideoId int64 `json:"video_id"` //视频id
	Cancel  int8  `json:"cancel"`   //是否点赞，0为点赞，1为取消赞

}

/*type Favourite struct {
	video
	user
}*/

// TableName 修改表名映射
func (f *Favourite) TableName() string {
	return consts.FavouriteTableName
}

// GetLikeUserIdList 根据videoId获取所有点赞userId
func GetLikeUserIdList(videoId int64) ([]int64, error) {
	return nil, nil
}

// GetLikeInfo 根据userId,videoId查询点赞信息
func GetLikeInfo(userId int64, videoId int64) (Favourite, error) {
	var favouriteInfo Favourite
	return favouriteInfo, nil
}

// UpdateLike 根据userId，videoId,actionType点赞或者取消赞
func UpdateLike(userId int64, videoId int64, actionType int32) error {
	return nil
}

// GetLikeVideoIdList 根据userId查询所属点赞全部videoId
func GetLikeVideoIdList(userId int64) ([]int64, error) {
	return nil, nil
}
