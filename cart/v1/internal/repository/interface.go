package repository

import "cart-service/v1/internal/constant"

type cartRepository interface {
	GetOrCreateCartByUserID(userID string) (*constant.Cart, error)
	AddItem(userID string, items []*constant.Item) error
	RemoveItem(userID string, cartID, itemID int) error
	RemoveCart(userID string) error
}
