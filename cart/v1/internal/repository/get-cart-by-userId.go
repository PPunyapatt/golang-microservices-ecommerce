package repository

import (
	"cart/v1/internal/constant"
)

func (repo *Repository) GetOrCreateCartByUserID(userID string) (*constant.Cart, error) {

	var cart constant.Cart
	args := []interface{}{userID}
	result := repo.GormDB.Where("user_id = ?", args...).First(&cart)
	if result.Error != nil {
		// if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		// }
		return nil, result.Error
	}

	return &cart, nil
}
