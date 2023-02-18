package db

import (
	"context"
	"log"
	"simple-tiktok/biz/model/tiktok"
	"time"

	"gorm.io/hints"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	UserId    uint      `gorm:"foreignkey:user(ID)"`
	User      *User
	Content   string
	VideoId   uint `gorm:"foreignkey:videos(ID)"`
}

func CreateComment(c context.Context, user_id uint, video_id uint, comment_text string) (*tiktok.Comment, error) {
	res := new(tiktok.Comment)
	res.User = new(tiktok.User)
	res.User = GetTiktokUser(c, user_id)
	res.Content = comment_text
	tmp := &Comment{UserId: user_id, Content: comment_text, VideoId: video_id}
	r := DB.Omit("DeletedAt").Create(&tmp)
	if r.Error != nil {
		return nil, r.Error
	}
	log.Println(tmp.CreatedAt)
	res.CreateDate = tmp.CreatedAt.Format("01-02")
	log.Println(res.CreateDate)
	res.ID = int64(tmp.ID)
	return res, nil
}

func DeleteComment(c context.Context, c_id uint) {
	DB.Delete(&Comment{}, c_id)
}

func GetComment(c context.Context, v_id uint, user_id uint) ([]*tiktok.Comment, error) {
	comments := make([]*Comment, 0)
	lf_list := make([]*lf_res, 0)
	DB.WithContext(c).Clauses(hints.UseIndex("idx_created_at")).Order("created_at desc").
		Where("video_id = ?", v_id).Preload("User").Find(&comments)
	query_fe := DB.WithContext(c).Table("follow").Select("COUNT(*) as follower_count, follow.user_id").Group("follow.user_id")
	query_f := DB.WithContext(c).Table("follow").Select("COUNT(*) as follow_count, follow.follower_id").Group("follow.follower_id")
	DB.WithContext(c).Table("comments").Order("comments.created_at desc").
		Select("follow.id as follow, q_f.follow_count, q_fe.follower_count").
		Joins("left join follow on follow.follower_id = ? AND comments.user_id=follow.user_id", user_id).
		Joins("left join (?) q_f on comments.user_id = q_f.follower_id", query_f).
		Joins("left join (?) q_fe on comments.user_id = q_fe.user_id", query_fe).
		Where("video_id = ?", v_id).Find(&lf_list)
	c_list := FormTiktokComment(comments, lf_list)
	return c_list, nil
}

func FormTiktokComment(comments []*Comment, lf_list []*lf_res) []*tiktok.Comment {
	c_list := make([]*tiktok.Comment, len(lf_list))
	for i := 0; i < len(lf_list); i++ {
		tmp := new(tiktok.Comment)
		tmp.Content = comments[i].Content
		tmp.CreateDate = comments[i].CreatedAt.Format("01-02")
		tmp.ID = int64(comments[i].ID)
		tmp.User = new(tiktok.User)
		tmp.User.Name = comments[i].User.Username
		tmp.User.ID = int64(comments[i].UserId)
		tmp.User.FollowCount = &lf_list[i].FollowCount
		tmp.User.FollowerCount = &lf_list[i].FollowerCount
		tmp.User.IsFollow = (lf_list[i].Follow > 0)
		c_list[i] = tmp
	}
	return c_list
}
