package repository

import (
	"errors"
	"fmt"
	"time"

	dao "chat-system/server/internal/infrastructure"
	userdomain "chat-system/server/internal/user/domain"

	"gorm.io/gorm"
)

type MySQLRepository struct{}

func NewMySQLRepository() MySQLRepository {
	return MySQLRepository{}
}

func (MySQLRepository) Create(user *userdomain.User) error {
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		var last userdomain.User
		account := "10000001"
		if err := tx.Table(userdomain.TableName).Order("account desc").First(&last).Error; err == nil {
			var number int
			if _, scanErr := fmt.Sscanf(last.Account, "%d", &number); scanErr == nil && number >= 10000001 {
				account = fmt.Sprintf("%08d", number+1)
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
		user.Account = account
		return tx.Table(userdomain.TableName).Create(user).Error
	})
}

func (MySQLRepository) FindByAccount(account string) (*userdomain.User, error) {
	var user userdomain.User
	if err := dao.DB.Table(userdomain.TableName).Where("account = ?", account).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (MySQLRepository) FindByID(userID uint) (*userdomain.User, error) {
	var user userdomain.User
	if err := dao.DB.Table(userdomain.TableName).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (MySQLRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return dao.DB.Table(userdomain.TableName).Where("id = ?", userID).Update("last_login_time", now).Error
}
