package db

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
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
/*func IsFavourite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	likeData := new(Favourite)
	//未查询到数据，返回未点赞；
	if result := DB.WithContext(ctx).Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likeData); result.RowsAffected == 0 {
		return false, errors.New("can't find this data")
	} //查询到数据，根据Cancel值判断是否点赞；
	if likeData.IsLike == consts.FavouriteAction {
		return true, nil
	} else {
		return false, nil
	}
}*/

// GetLikeUserIdList 根据videoId获取所有点赞userId
func GetLikeUserIdList(ctx context.Context, videoId int64) ([]int64, error) {
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

}

/*// GetLikeInfo 根据userId,videoId查询点赞信息
func GetLikeInfo(userId int64, videoId int64) (Favourite, error) {
	var favouriteInfo Favourite
	return favouriteInfo, nil
}*/

// FavouriteAction 根据userId，videoId点赞或者取消赞
func FavouriteAction(ctx context.Context, userId int64, videoId int64, action_type int32) error {
	likeData := new(Favourite)
	//先查询是否有这条数据。
	result := DB.WithContext(ctx).Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likeData)
	//点赞行为,否则取消赞
	if action_type == consts.LikeAction {
		//没查到这条数据，则新建这条点赞数据,否则更新即可;
		likeData.UserId = userId
		likeData.VideoId = videoId
		likeData.IsLike = consts.FavouriteAction
		if result1 := DB.WithContext(ctx).Table("likes").Create(&likeData); result1.RowsAffected == 0 {
			return errors.New("insert data fail")
		} else {
			//查询到数据，更新点赞信息
			if result2 := DB.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).
				Update("cancel", consts.FavouriteAction); result2.RowsAffected == 0 {
				return errors.New("update data fail")
			}
		}
	} else {
		//只有当前是点赞状态才能取消点赞这个行为，如果查询不到数据则返回错误；
		if result.RowsAffected == 0 {
			return errors.New("can't find this data")
		} else {
			//若已经点赞过，取消赞	 数据库中"cancel"改为0未点赞的状态
			if result3 := DB.WithContext(ctx).Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).
				Update("cancel", 0); result3.RowsAffected == 0 {
				return errors.New("update data fail")
			}
		}
	}
	return nil
}

// GetFavouriteList  根据userId查询like表中点赞的全部videoId
func GetFavouriteList(ctx context.Context, userId int64) ([]Video, error) {
	results := make([]Video, 0)
	err := DB.WithContext(ctx).Model(&Favourite{}).Select("user_id").Where("user_id in ? and video_id = ?", userId).Find(&results).Error
	return results, err
}
