package repository

import (
	"errors"

	"chat-system/server/internal/common"
	frienddomain "chat-system/server/internal/friend/domain"
	dao "chat-system/server/internal/infrastructure"
	userdomain "chat-system/server/internal/user/domain"

	"gorm.io/gorm"
)

type MySQLRepository struct{}

func NewMySQLRepository() MySQLRepository {
	return MySQLRepository{}
}

func (MySQLRepository) Request(userID uint, friendID uint) error {
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		var existing frienddomain.Friend
		err := tx.Table(frienddomain.TableName).Where("user_id = ? AND friend_id = ?", userID, friendID).First(&existing).Error
		if err == nil {
			if existing.Status == common.FriendStatusAccepted {
				return errors.New("already friends")
			}
			return errors.New("friend request already sent")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var reverse frienddomain.Friend
		err = tx.Table(frienddomain.TableName).Where("user_id = ? AND friend_id = ?", friendID, userID).First(&reverse).Error
		if err == nil {
			if reverse.Status == common.FriendStatusAccepted {
				return errors.New("already friends")
			}
			return acceptRequestTx(tx, reverse.ID, userID)
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		request := frienddomain.Friend{UserID: userID, FriendID: friendID, Status: common.FriendStatusPending}
		return tx.Table(frienddomain.TableName).Create(&request).Error
	})
}

func (MySQLRepository) AcceptRequest(requestID uint, userID uint) error {
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		return acceptRequestTx(tx, requestID, userID)
	})
}

func (MySQLRepository) ExistsAccepted(userID uint, friendID uint) (bool, error) {
	var count int64
	err := dao.DB.Table(frienddomain.TableName).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, common.FriendStatusAccepted).
		Count(&count).Error
	return count > 0, err
}

func acceptRequestTx(tx *gorm.DB, requestID uint, userID uint) error {
	var request frienddomain.Friend
	if err := tx.Table(frienddomain.TableName).Where("id = ? AND friend_id = ?", requestID, userID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("friend request not found")
		}
		return err
	}
	if request.Status == common.FriendStatusAccepted {
		return errors.New("already friends")
	}

	if err := tx.Table(frienddomain.TableName).Where("id = ?", request.ID).Update("status", common.FriendStatusAccepted).Error; err != nil {
		return err
	}

	reverse := frienddomain.Friend{UserID: userID, FriendID: request.UserID, Status: common.FriendStatusAccepted}
	return tx.Table(frienddomain.TableName).Where("user_id = ? AND friend_id = ?", reverse.UserID, reverse.FriendID).
		Assign(frienddomain.Friend{Status: common.FriendStatusAccepted}).
		FirstOrCreate(&reverse).Error
}

func (MySQLRepository) ListAccepted(userID uint) ([]frienddomain.Friend, error) {
	var friends []frienddomain.Friend
	if err := dao.DB.Table(frienddomain.TableName).Where("user_id = ? AND status = ?", userID, common.FriendStatusAccepted).Find(&friends).Error; err != nil {
		return nil, err
	}
	return withUsers(friends)
}

func (MySQLRepository) ListPendingRequests(userID uint) ([]frienddomain.Friend, error) {
	var friends []frienddomain.Friend
	if err := dao.DB.Table(frienddomain.TableName).Where("friend_id = ? AND status = ?", userID, common.FriendStatusPending).Find(&friends).Error; err != nil {
		return nil, err
	}
	return withUsers(friends)
}

func withUsers(friends []frienddomain.Friend) ([]frienddomain.Friend, error) {
	for index := range friends {
		if err := dao.DB.Table(userdomain.TableName).Where("id = ?", friends[index].UserID).First(&friends[index].User).Error; err != nil {
			return nil, err
		}
		if err := dao.DB.Table(userdomain.TableName).Where("id = ?", friends[index].FriendID).First(&friends[index].Friend).Error; err != nil {
			return nil, err
		}
	}
	return friends, nil
}
