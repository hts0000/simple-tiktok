package db

import (
	"context"
	"errors"
	"simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string `json:"username"`
	Password           string `json:"password"`
	FollowCount        uint   `json:"follow_count"`
	FollowerCount      uint   `json:"follower_count"`
	AvatarURL          string `json:"avatar_url" gorm:"default:https://simple-tiktok-1300912551.cos.ap-guangzhou.myqcloud.com/avatar.jpg"`
	BackgroundImageURL string `json:"background_image_url" gorm:"default:https://simple-tiktok-1300912551.cos.ap-guangzhou.myqcloud.com/background_image.jpg"`
	Signature          string `json:"signature" gorm:"default:我是一只抖小萌😘💗💓"`
	TotalFavorited     uint   `json:"total_favorited"`
	WorkCount          uint   `json:"work_count"`
	FavoriteCount      uint   `json:"favorite_count"`
}

func (u *User) TableName() string {
	return consts.UserTableName
}

func GetUser(ctx context.Context, userID int64) (*User, error) {
	user := User{}
	err := DB.WithContext(ctx).Where("id = ?", userID).Take(&user).Error
	return &user, err
}

// MGetUsers multiple get list of user info
func MGetUsers(c context.Context, userIDs []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	result := DB.WithContext(c).Where("id in ?", userIDs).Find(&res)
	err := result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) || err == nil {
		return res, nil
	}
	return res, err
}

// CreateUser create user info
func CreateUser(c context.Context, username, password string) (uint, error) {
	user := &User{
		Username: username,
		Password: password,
	}
	err := DB.WithContext(c).Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// QueryUser query list of user info
func QueryUser(c context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)
	if err := DB.WithContext(c).Where("username = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
