package repository

import (
	"cart/v1/internal/constant"
	"cart/v1/proto/cart"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func (repo *Repository) GetItemsByUserID(userID string, pagination *constant.Pagination) ([]*cart.CartItem, error) {
	query := `
		SELECT 
			ci.product_id,
			ci.product_name,
			ci.quantity,
			ci.price,
			ci.image_url,
			ci.store_id
		FROM carts c
		LEFT JOIN cart_items ci
			ON c.id = ci.cart_id AND 
			c.user_id = $1
	`

	args := []interface{}{userID}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM (%s)`, query)
	if err := repo.PostgresDB.QueryRowx(
		countQuery,
		args...,
	).Scan(&pagination.TotalCount); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	args = append(args, pagination.Limit, pagination.Offset)
	query += fmt.Sprintf(` LIMIT $2 OFFSET $3`)

	items := []*constant.Item{}
	if err := repo.PostgresDB.Select(&items, query, args...); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	cartItems := []*cart.CartItem{}
	for _, item := range items {
		cartItems = append(cartItems, &cart.CartItem{
			ProductId:   int32(item.ProductID),
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			Price:       item.Price,
			ImageUrl:    item.ImageURL,
			// StoreID:     int32(item.StoreID),
		})
	}

	return cartItems, nil
}
