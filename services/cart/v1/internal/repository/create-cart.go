package repository

import (
	"cart/v1/internal/constant"
)

func (repo *Repository) CreateCart(userID string) error {
	cart := &constant.Cart{
		UserID: userID,
	}

	result := repo.GormDB.Create(cart)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
