package db

import (
	"context"
	"fmt"
	tiktok "simple-tiktok/biz/model/tiktok"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title          string
	AuthorId       uint `gorm:"foreignkey:user(ID)"`
	Author         User
	Play_url       string
	Cover_url      string
	Favorite_count int64
	Comment_count  int64
}

type lf_res struct {
	VideoId       uint
	AuthorId      uint
	Cancel        int8
	Follow        uint
	FollowCount   int64
	FollowerCount int64
}

func CreateVideoAndGetId(c context.Context, title string, tp string, user_id int64) (int64, error) {
	var v Video
	au, err := GetUser(c, user_id)
	if err != nil {
		return -1, err
	}
	p_url := "http://172.20.15.52:8080/file/?location=video/"
	c_url := "http://172.20.15.52:8080/file/?location=img/"
	v = Video{Title: title, AuthorId: au.ID, Play_url: p_url, Cover_url: c_url}
	res := DB.WithContext(c).Select("Title", "AuthorId", "Play_url", "Cover_url").Create(&v)
	if res.Error != nil {
		return -1, res.Error
	}
	v.Play_url = fmt.Sprintf("http://172.20.15.52:8080/file/?location=video/%s.%s", strconv.FormatInt(int64(v.ID), 10), tp)
	v.Cover_url = fmt.Sprintf("http://172.20.15.52:8080/file/?location=img/%s.jpg", strconv.FormatInt(int64(v.ID), 10))
	DB.WithContext(c).Model(&Video{}).Where("ID = ?", v.ID).Updates(v)
	return int64(v.ID), nil
}

func GetFeedVideo(c context.Context, tm time.Time, user_id int64) ([]*tiktok.Video, time.Time, error) {
	videos := make([]*Video, 0)
	lf_list := make([]*lf_res, 0)
	var err error
	DB.WithContext(c).Order("created_at desc").Where("created_at > ?", tm).Preload("Author").Limit(30).Find(&videos)
	query_f := DB.WithContext(c).Table("follow").Select("COUNT(*) as follow_count, follow.user_id").Group("follow.user_id")
	query_fe := DB.WithContext(c).Table("follow").Select("COUNT(*) as follower_count, follow.follower_id").Group("follow.follower_id")
	DB.WithContext(c).Table("videos").Order("videos.created_at desc").Select("videos.id as video_id, videos.author_id, likes.cancel, follow.id as follow, q_f.follow_count, q_fe.follower_count").
		Joins("left join likes on videos.ID=likes.video_id AND videos.author_id=likes.user_id").Joins("left join follow on follow.user_id = ? AND videos.author_id=follow.follower_id", uint(user_id)).
		Joins("left join (?) q_f on videos.author_id = q_f.user_id", query_f).Joins("left join (?) q_fe on videos.author_id = q_fe.follower_id", query_fe).
		Where("videos.created_at > ?", tm).Limit(30).Find(&lf_list)
	v_list := make([]*tiktok.Video, len(videos))
	for i := 0; i < len(v_list); i++ {
		tmp := new(tiktok.Video)
		tmp.Author = new(tiktok.User)
		tmp.Author.ID = int64(videos[i].Author.ID)
		tmp.Author.Name = videos[i].Author.Username
		tmp.Author.IsFollow = (lf_list[i].Follow > 0)
		tmp.Author.FollowCount = &(lf_list[i].FollowCount)
		tmp.Author.FollowerCount = &(lf_list[i].FollowerCount)
		tmp.CommentCount = videos[i].Comment_count
		tmp.FavoriteCount = videos[i].Favorite_count
		tmp.ID = int64(videos[i].ID)
		tmp.Title = videos[i].Title
		tmp.PlayURL = videos[i].Play_url
		tmp.CoverURL = videos[i].Cover_url
		tmp.IsFavorite = (lf_list[i].Cancel > 0)
		v_list[i] = tmp
	}
	var latest_time time.Time
	if len(videos) > 0 {
		latest_time = videos[len(videos)-1].CreatedAt
	}
	return v_list, latest_time, err
}
