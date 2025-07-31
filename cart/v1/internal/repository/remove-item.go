package repository

import "cart/v1/internal/constant"

func (repo *Repository) RemoveItem(userID string, cartID, itemID int) error {
	err := repo.GormDB.Where("id = ?", itemID).Delete(&constant.Item{}).Error
	if err != nil {
		return err
	}

	var count int64
	err = repo.GormDB.Model(&constant.Item{}).Where("cart_id = ?", cartID).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		err = repo.RemoveCart(userID)
		if err != nil {
			return err
		}
	}

	return nil
}
