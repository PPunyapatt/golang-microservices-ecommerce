package repository

import (
	"cart-service/v1/internal/constant"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (repo *Repository) GetOrCreateCartByUserID(userID string) (*constant.Cart, error) {

	var cart constant.Cart
	args := []interface{}{userID}
	result := repo.GormDB.Where("user_id = ?", args...).First(&cart)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			cart = constant.Cart{
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			resultCreate := repo.GormDB.Create(&cart)
			if resultCreate.Error != nil {
				return nil, result.Error
			}
		}
	}

	return &cart, nil
}
