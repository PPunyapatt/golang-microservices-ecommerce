package repository

import (
	"cart/v1/internal/constant"
	"log"

	"github.com/pkg/errors"
)

func (repo *Repository) RemoveItem(userID string, cartID, itemID int) error {
	query := `
		DELETE FROM cart_items ci
		USING carts c
		WHERE 
			c.id = ci.cart_id AND
			c.user_id = $1 AND
			ci.cart_id = $2 AND
			ci.id = $3
	`

	args := []interface{}{userID, cartID, itemID}
	result, err := repo.PostgresDB.Exec(query, args...)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected, item not found or already removed")
	}

	var count int64
	if err := repo.GormDB.Model(&constant.Item{}).
		Where("cart_id = ?", cartID).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := repo.RemoveCart(userID); err != nil {
			return err
		}
	}

	return nil
}
