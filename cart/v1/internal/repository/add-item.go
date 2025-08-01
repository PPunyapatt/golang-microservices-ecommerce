package repository

import (
	"cart/v1/internal/constant"
)

func (repo *Repository) AddItem(userID string, items []*constant.Item) error {

	cart, err := repo.GetOrCreateCartByUserID(userID)
	if err != nil {
		return err
	}

	for _, item := range items {
		item.CartID = cart.CartID
	}
	// godump.Dump(items)
	if err := repo.GormDB.Omit("updated_at").Create(&items); err != nil {
		return err.Error
	}
	return nil
}
