package db

import (
	"context"
	"errors"
	"simple-tiktok/pkg/consts"
	"time"

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
	message := Message{}
	result := DB.WithContext(ctx).Where("sender_id IN ? AND receiver_id IN ?", []uint{u1, u2}, []uint{u1, u2}).Last(&message)
	err := result.Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return message.Message, nil
	}
	return "", err
}

// 获取u1和u2之间时间在timestamp之后的消息列表
func GetMessages(ctx context.Context, u1, u2 uint, timestamp int64) ([]*Message, error) {
	messages := make([]*Message, 0)
	lastTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	result := DB.WithContext(ctx).
		Where("sender_id IN ? AND receiver_id IN ?", []uint{u1, u2}, []uint{u1, u2}).
		Where("updated_at > ?", lastTime).
		Find(&messages)
	return messages, result.Error
}

// u1给u2发生一条消息
func CreateMessage(ctx context.Context, u1, u2 uint, msg string) (uint, error) {
	message := Message{
		SenderID:   u1,
		ReceiverID: u2,
		Message:    msg,
	}
	err := DB.WithContext(ctx).Create(&message).Error
	return message.ID, err
}
