package repository

import "cart/v1/internal/constant"

func (repo *Repository) AddItem(userID string, items []*constant.Item) error {

	cart, err := repo.GetOrCreateCartByUserID(userID)
	if err != nil {
		return err
	}

	for _, item := range items {
		item.CartID = cart.CartID
	}

	if err := repo.GormDB.Create(&items); err != nil {
		return err.Error
	}
	return nil
}
