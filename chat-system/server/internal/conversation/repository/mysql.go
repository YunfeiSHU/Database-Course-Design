package repository

import (
	conversationdomain "chat-system/server/internal/conversation/domain"
	dao "chat-system/server/internal/infrastructure"
)

type MySQLRepository struct{}

func NewMySQLRepository() MySQLRepository {
	return MySQLRepository{}
}

func (MySQLRepository) Upsert(conversation conversationdomain.Conversation) error {
	return dao.DB.Table(conversationdomain.TableName).
		Where("user_id = ? AND peer_id = ?", conversation.UserID, conversation.PeerID).
		Assign(map[string]interface{}{
			"last_message_id": conversation.LastMessageID,
			"status":          conversation.Status,
			"update_time":     conversation.UpdateTime,
		}).
		FirstOrCreate(&conversation).Error
}

func (MySQLRepository) ListByUser(userID uint) ([]conversationdomain.Conversation, error) {
	var conversations []conversationdomain.Conversation
	err := dao.DB.Table(conversationdomain.TableName).
		Where("user_id = ?", userID).
		Order("update_time desc").
		Find(&conversations).Error
	return conversations, err
}
