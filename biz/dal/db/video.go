package db

import (
	"context"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	User_id int64  `json:"user_id"`
}

func CreateVideoAndGetId(c context.Context, title string, tp string, user_id int64) (int64, error) {
	tp = "." + tp
	v := Video{Title: title, Type: tp, User_id: user_id}
	result := DB.WithContext(c).Create(&v)
	if result.Error != nil {
		return -1, result.Error
	}
	return v.ID, nil
}
