package db

import (
	"context"
	"fmt"
	"log"
	tiktok "simple-tiktok/biz/model/tiktok"
	"strconv"
	"strings"
	"time"

	"gorm.io/hints"

	"github.com/go-redis/redis"
)

type Video struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	Title     string
	AuthorId  uint `gorm:"foreignkey:user(ID)"`
	Author    *User
	Play_url  string
	Cover_url string
}

type lf_res struct {
	Cancel        int8
	Follow        uint
	FollowCount   int64
	FollowerCount int64
	CommentCount  int64
	FavoriteCount int64
}

func CreateVideoAndGetId(c context.Context, title string, tp string, user_id int64) (string, error) {
	var v Video
	fileName := time.Now().Format("2006-01-02 15:04:05")
	fileName = strings.ReplaceAll(fileName, "-", "")
	fileName = strings.ReplaceAll(fileName, " ", "")
	fileName = strings.ReplaceAll(fileName, ":", "")
	p_url := fmt.Sprintf("http://172.20.15.52:8080/file/?location=video/%s.%s", fileName, tp)
	c_url := fmt.Sprintf("http://172.20.15.52:8080/file/?location=img/%s.jpg", fileName)
	v = Video{Title: title, AuthorId: uint(user_id), Play_url: p_url, Cover_url: c_url}
	res := DB.WithContext(c).Select("Title", "AuthorId", "Play_url", "Cover_url").Create(&v)
	if res.Error != nil {
		return "", res.Error
	}
	RD.ZAdd("TimeLine", redis.Z{Score: float64(v.CreatedAt.Unix()), Member: strconv.FormatInt(int64(v.ID), 10)})
	key_v := "video_" + strconv.FormatInt(int64(v.ID), 10) + "_info"
	_, err := RD.HMSet(key_v, map[string]interface{}{
		"Title":         v.Title,
		"PlayURL":       v.Play_url,
		"CoverURL":      v.Cover_url,
		"CommentCount":  0,
		"FavoriteCount": 0,
		"AuthorId":      user_id,
	}).Result()
	RD.Expire(key_v, 5*time.Second)
	if err != nil {
		log.Printf("Insert into redis failed! err: %s.\n", err.Error())
		return "", err
	}
	return fileName, nil
}

func GetFeedVideo(c context.Context, max string, user_id int64) ([]*tiktok.Video, int64, error) {
	min := "-inf"
	id_list, err := RD.ZRevRangeByScoreWithScores("TimeLine", redis.ZRangeBy{Min: min, Max: max, Offset: 0, Count: 30}).Result()
	if err != nil {
		return nil, time.Now().Unix(), err
	}
	v_list := make([]*tiktok.Video, len(id_list))
	for i := 0; i < len(v_list); i++ {
		v_list[i] = new(tiktok.Video)
		log.Println(id_list[i].Member)
		au_id, err := GetTiktokVideo(c, v_list[i], id_list[i].Member.(string), uint(user_id))
		if err != nil {
			log.Printf("Get video infomation failed! err: %s.\n", err.Error())
			return nil, time.Now().Unix(), err
		}
		log.Println(au_id)
		v_list[i].Author, err = NewGetTiktokUser(c, au_id, user_id)
		if err != nil {
			log.Printf("Get author infomation failed! err: %s.\n", err.Error())
			return nil, time.Now().Unix(), err
		}
		log.Println(*v_list[i].Author)
		log.Println(*v_list[i])
	}
	return v_list, int64(id_list[len(id_list)-1].Score), err
}

func GetPublishList(c context.Context, user_id uint) ([]*tiktok.Video, error) {
	var err error
	videos := make([]*Video, 0)
	lf_list := make([]*lf_res, 0)
	DB.WithContext(c).Where("author_id = ?", user_id).Find(&videos)
	query_c := DB.WithContext(c).Table("comments").Select("COUNT(*) as comment_count, comments.video_id").Group("comments.video_id")
	query_l := DB.WithContext(c).Table("likes").Select("Count(*) as favorite_count, likes.video_id").Group("likes.video_id")
	DB.WithContext(c).Table("videos").Where("videos.author_id = ?", user_id).
		Select("videos.id as video_id, likes.cancel, q_c.comment_count, q_l.favorite_count").
		Joins("left join likes on videos.ID=likes.video_id AND videos.author_id=likes.user_id").
		Joins("left join (?) q_c on videos.ID = q_c.video_id", query_c).
		Joins("left join (?) q_l on videos.ID = q_l.video_id", query_l).
		Find(&lf_list)
	author := GetTiktokUser(c, user_id)
	log.Println(author)
	v_list := PublishListTiktokVideo(videos, lf_list, *author)
	return v_list, err
}

func PublishListTiktokVideo(videos []*Video, lf_list []*lf_res, author tiktok.User) []*tiktok.Video {
	v_list := make([]*tiktok.Video, len(lf_list))
	for i := 0; i < len(v_list); i++ {
		tmp := new(tiktok.Video)
		tmp.Author = new(tiktok.User)
		tmp.Author.ID = author.ID
		tmp.Author.Name = author.Name
		tmp.Author.IsFollow = author.IsFollow
		tmp.Author.FollowCount = author.FollowCount
		tmp.Author.FollowerCount = author.FollowerCount
		tmp.CommentCount = lf_list[i].CommentCount
		tmp.FavoriteCount = lf_list[i].FavoriteCount
		tmp.ID = int64(videos[i].ID)
		tmp.Title = videos[i].Title
		tmp.PlayURL = videos[i].Play_url
		tmp.CoverURL = videos[i].Cover_url
		tmp.IsFavorite = (lf_list[i].Cancel > 0)
		v_list[i] = tmp
	}
	return v_list
}

func GetTiktokUser(c context.Context, user_id uint) *tiktok.User {
	author := new(tiktok.User)
	author.FollowCount = new(int64)
	author.FollowerCount = new(int64)
	author.ID = int64(user_id)
	DB.WithContext(c).Table("user").Select("user.id, user.username as name").Where("id = ?", user_id).Find(&author)
	DB.WithContext(c).Table("follow").Where("follow.follower_id = ?", user_id).Count(author.FollowCount)
	DB.WithContext(c).Table("follow").Where("follow.user_id = ?", user_id).Count(author.FollowerCount)
	DB.WithContext(c).Table("follow").Select("follow.id as is_follow").Where("follow.follower_id = ? AND follow.user_id = ?", user_id, user_id).Find(&author)
	return author
}

func NewGetTiktokUser(c context.Context, au_id int64, user_id int64) (*tiktok.User, error) {
	key_u := "user_" + strconv.FormatInt(au_id, 10) + "_info"
	author := new(tiktok.User)
	author.ID = au_id
	res, err := RD.HGetAll(key_u).Result()
	log.Printf("whether user %d find redis: %d", au_id, len(res))
	if err == nil {
		if len(res) == 0 {
			DB.WithContext(c).Table("user").Select(`user.username as name, user.avatar_url as avatar,
				user.background_image_url as background_image, user.signature, user.total_favorited, 
				user.work_count, user.favorite_count`).Where("id = ?", uint(author.ID)).Find(&author)
			log.Println(*author)
			author.FollowCount = new(int64)
			DB.WithContext(c).Clauses(hints.UseIndex("idx_follower_id")).Table("follow").Where("follow.follower_id = ?", uint(author.ID)).Count(author.FollowCount)
			author.FollowerCount = new(int64)
			DB.WithContext(c).Clauses(hints.UseIndex("idx_user_id")).Table("follow").Where("follow.user_id = ?", uint(author.ID)).Count(author.FollowerCount)
			_, err = RD.HMSet(key_u, map[string]interface{}{
				"Name":            author.Name,
				"FollowCount":     *author.FollowCount,
				"FollowerCount":   *author.FollowerCount,
				"Avatar":          *author.Avatar,
				"BackgroundImage": *author.BackgroundImage,
				"Signature":       *author.Signature,
				"TotalFavorite":   *author.TotalFavorited,
				"WorkCount":       *author.WorkCount,
				"FavoriteCount":   *author.FavoriteCount,
			}).Result()
			RD.Expire(key_u, 10*time.Second)
			if err != nil {
				log.Printf("Insert Redis Failed, error:%s!\n", err.Error())
				return nil, err
			}
		} else {
			author.Name = res["Name"]
			follow_count, _ := strconv.ParseInt(res["FollowCount"], 10, 64)
			follower_count, _ := strconv.ParseInt(res["FollowerCount"], 10, 64)
			author.FollowCount = &follow_count
			author.FollowerCount = &follower_count
			Avatar := res["Avatar"]
			author.Avatar = &Avatar
			Bgi := res["BackgroundImage"]
			author.BackgroundImage = &Bgi
			Sig := res["Signature"]
			author.Signature = &Sig
			WorkCount, _ := strconv.ParseInt(res["WorkCount"], 10, 64)
			author.WorkCount = &WorkCount
			FavoriteCount, _ := strconv.ParseInt(res["FavoriteCount"], 10, 64)
			author.FavoriteCount = &FavoriteCount
			Total, _ := strconv.ParseInt(res["TotalFavorite"], 10, 64)
			author.TotalFavorited = &Total
		}
	} else {
		log.Printf("Select Redis Failed, error:%s!\n", err.Error())
		return nil, err
	}
	//后续建好整个redis后，也改用缓存
	result := map[string]interface{}{}
	DB.WithContext(c).Table("follow").Select("follow.id as is_follow").Omit("Author").Where("follow.follower_id = ? AND follow.user_id = ?", uint(au_id), uint(user_id)).Find(&result)
	log.Println(result)
	author.IsFollow = (len(result) != 0)
	return author, nil
}

// 查找单个视频
func GetVideo(c context.Context, vid uint) (*Video, error) {
	video := Video{
		ID: vid,
	}
	err := DB.WithContext(c).Take(&video).Error
	return &video, err
}

func GetTiktokVideo(c context.Context, vd *tiktok.Video, v_id string, user_id uint) (int64, error) {
	key_v := "video_" + v_id + "_info"
	vd.ID, _ = strconv.ParseInt(v_id, 10, 64)
	res, err := RD.HGetAll(key_v).Result()
	log.Printf("whether video %d find redis: %d", vd.ID, len(res))
	var au_id int64
	if err == nil {
		if len(res) == 0 {
			DB.WithContext(c).Clauses(hints.UseIndex("idx_video")).Table("comments").Where("comments.video_id = ?", uint(vd.ID)).Count(&vd.CommentCount)
			DB.WithContext(c).Clauses(hints.UseIndex("videoIdx")).Table("likes").Where("likes.video_id = ?", uint(vd.ID)).Count(&vd.FavoriteCount)
			video := Video{
				ID: uint(vd.ID),
			}
			err = DB.WithContext(c).Take(&video).Error
			if err != nil {
				log.Printf("Select Failed, error:%s!\n", err.Error())
				return -1, err
			}
			vd.PlayURL = video.Play_url
			vd.CoverURL = video.Cover_url
			vd.Title = video.Title
			au_id = int64(video.AuthorId)
			_, err = RD.HMSet(key_v, map[string]interface{}{
				"Title":         vd.Title,
				"PlayURL":       vd.PlayURL,
				"CoverURL":      vd.CoverURL,
				"CommentCount":  vd.CommentCount,
				"FavoriteCount": vd.FavoriteCount,
				"AuthorId":      au_id,
			}).Result()
			RD.Expire(key_v, 5*time.Second)
			if err != nil {
				log.Printf("Insert Redis Failed, error:%s!\n", err.Error())
				return -1, err
			}
		} else {
			vd.CommentCount, _ = strconv.ParseInt(res["CommentCount"], 10, 64)
			vd.FavoriteCount, _ = strconv.ParseInt(res["FavoriteCount"], 10, 64)
			vd.Title = res["Title"]
			vd.CoverURL = res["CoverURL"]
			vd.PlayURL = res["PlayURL"]
			au_id, _ = strconv.ParseInt(res["AuthorId"], 10, 64)
		}
	} else {
		log.Printf("Select Redis Failed, error:%s!\n", err.Error())
		return -1, err
	}
	//后续建好整个redis后，也改用缓存
	result := map[string]interface{}{}
	DB.WithContext(c).Table("likes").Select("likes.id as is_favorite").
		Where("likes.user_id = ? AND likes.video_id = ?", user_id, vd.ID).Find(&result)
	vd.IsFavorite = (len(result) != 0)
	return au_id, nil
}
