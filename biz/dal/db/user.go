package db

import (
	"context"
	"errors"
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
