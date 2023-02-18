package db

import (
	"context"
	"errors"
	"simple-tiktok/pkg/consts"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Message    string `json:"message"`
}

func (m *Message) TableName() string {
	return consts.MessageTableName
}

// 获取u1和u2之间的最新消息
func GetMessage(ctx context.Context, u1, u2 uint) (string, error) {
	message1 := Message{}
	message2 := Message{}
	result1 := DB.WithContext(ctx).Where("sender_id = ? AND receiver_id = ?", u1, u2).Last(&message1)
	result2 := DB.WithContext(ctx).Where("sender_id = ? AND receiver_id = ?", u2, u1).Last(&message2)
	err1, err2 := result1.Error, result2.Error
	if err1 == nil && err2 == nil {
		if message1.ID > message2.ID {
			return message1.Message, nil
		}
		return message2.Message, nil
	}
	if err1 == nil {
		return message1.Message, nil
	} else if err2 == nil {
		return message2.Message, nil
	}
	if errors.Is(err1, gorm.ErrRecordNotFound) && errors.Is(err2, gorm.ErrRecordNotFound) {
		return "", nil
	} else if !errors.Is(err1, gorm.ErrRecordNotFound) {
		return "", err1
	}
	return "", err2
}
