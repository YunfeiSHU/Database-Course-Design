package repository

import (
	dao "chat-system/server/internal/infrastructure"
	messagedomain "chat-system/server/internal/message/domain"
)

type MySQLRepository struct{}

func NewMySQLRepository() MySQLRepository {
	return MySQLRepository{}
}

func (MySQLRepository) Save(message *messagedomain.Message) error {
	return dao.DB.Table(messagedomain.TableName).Create(message).Error
}

func (MySQLRepository) FindByID(messageID uint) (*messagedomain.Message, error) {
	var message messagedomain.Message
	if err := dao.DB.Table(messagedomain.TableName).Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (MySQLRepository) UpdateStatus(messageID uint, status string) error {
	return dao.DB.Table(messagedomain.TableName).Where("id = ?", messageID).Update("status", status).Error
}

func (MySQLRepository) UpdateContentAndStatus(messageID uint, content string, status string) error {
	return dao.DB.Table(messagedomain.TableName).Where("id = ?", messageID).Updates(map[string]interface{}{
		"content": content,
		"status":  status,
	}).Error
}

func (MySQLRepository) List(userID uint, friendID uint, limit int) ([]messagedomain.Message, error) {
	var messages []messagedomain.Message
	err := dao.DB.Table(messagedomain.TableName).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, friendID, friendID, userID).
		Order("send_time desc").
		Limit(limit).
		Find(&messages).Error
	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}
	return messages, err
}
