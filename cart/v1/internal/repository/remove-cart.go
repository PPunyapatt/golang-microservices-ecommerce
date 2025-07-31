package repository

import "cart/v1/internal/constant"

func (repo *Repository) RemoveCart(userID string) error {
	err := repo.GormDB.Where("user_id = ?", userID).Delete(&constant.Cart{}).Error
	if err != nil {
		return err
	}

	return nil
}
