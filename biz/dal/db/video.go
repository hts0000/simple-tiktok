package db

import (
	"context"
	"fmt"
	"log"
	tiktok "simple-tiktok/biz/model/tiktok"
	"strings"
	"time"

	"gorm.io/hints"
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
	return fileName, nil
}

func GetFeedVideo(c context.Context, tm time.Time, user_id int64) ([]*tiktok.Video, time.Time, error) {
	videos := make([]*Video, 0)
	lf_list := make([]*lf_res, 0)
	var err error
	DB.WithContext(c).Clauses(hints.UseIndex("idx_videos_created_at")).Order("created_at desc").Where("created_at > ?", tm).Preload("Author").Limit(30).Find(&videos)
	query_fe := DB.WithContext(c).Table("follow").Select("COUNT(*) as follower_count, follow.user_id").Group("follow.user_id")
	query_f := DB.WithContext(c).Table("follow").Select("COUNT(*) as follow_count, follow.follower_id").Group("follow.follower_id")
	query_c := DB.WithContext(c).Table("comments").Select("COUNT(*) as comment_count, comments.video_id").Group("comments.video_id")
	query_l := DB.WithContext(c).Table("likes").Select("Count(*) as favorite_count, likes.video_id").Group("likes.video_id")
	DB.WithContext(c).Table("videos").Order("videos.created_at desc").
		Select("likes.cancel, follow.id as follow, q_f.follow_count, q_fe.follower_count, q_c.comment_count, q_l.favorite_count").
		Joins("left join likes on videos.ID=likes.video_id AND videos.author_id=likes.user_id").
		Joins("left join follow on follow.follower_id = ? AND videos.author_id=follow.user_id", uint(user_id)).
		Joins("left join (?) q_f on videos.author_id = q_f.follower_id", query_f).
		Joins("left join (?) q_fe on videos.author_id = q_fe.user_id", query_fe).
		Joins("left join (?) q_c on videos.ID = q_c.video_id", query_c).
		Joins("left join (?) q_l on videos.ID = q_l.video_id", query_l).
		Where("videos.created_at > ?", tm).Limit(30).Find(&lf_list)
	v_list := FeedTiktokVideo(videos, lf_list)
	var latest_time time.Time
	if len(videos) > 0 {
		latest_time = videos[len(videos)-1].CreatedAt
	}
	return v_list, latest_time, err
}

func FeedTiktokVideo(videos []*Video, lf_list []*lf_res) []*tiktok.Video {
	v_list := make([]*tiktok.Video, len(lf_list))
	for i := 0; i < len(v_list); i++ {
		tmp := new(tiktok.Video)
		tmp.Author = new(tiktok.User)
		tmp.Author.ID = int64(videos[i].Author.ID)
		tmp.Author.Name = videos[i].Author.Username
		tmp.Author.IsFollow = (lf_list[i].Follow > 0)
		tmp.Author.FollowCount = &(lf_list[i].FollowCount)
		tmp.Author.FollowerCount = &(lf_list[i].FollowerCount)
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
	log.Println(v_list[0])
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
