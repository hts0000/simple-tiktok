package db

import (
	"context"
	"simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) TableName() string {
	return consts.UserTableName
}

// MGetUsers multiple get list of user info
func MGetUsers(c context.Context, userIDs []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := DB.WithContext(c).Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateUser create user info
func CreateUser(c context.Context, user *User) (uint, error) {
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
