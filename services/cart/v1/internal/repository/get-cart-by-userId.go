package repository

import (
	"cart/v1/internal/constant"
)

func (repo *Repository) GetOrCreateCartByUserID(userID string) (*constant.Cart, error) {

	cart := constant.Cart{
		UserID: userID,
	}
	args := []interface{}{userID}
	result := repo.GormDB.Where("user_id = ?", args...).Omit("updated_at").FirstOrCreate(&cart)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cart, nil
}
