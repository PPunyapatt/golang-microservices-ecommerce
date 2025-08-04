package repository

import "context"

// func (repo *orderRepository) GetOrder(ctx context.Context, orderID int32) error {
// 	return nil
// }

// func (repo *orderRepository) ListOrder(ctx context.Context, status *string) error {
// 	return nil
// }

func (repo *orderRepository) AddOrderItems(ctx context.Context) error {
	return nil
}

func (repo *orderRepository) AddOrder(ctx context.Context) error {
	return nil
}

func (repo *orderRepository) ChangeStatus(ctx context.Context, orderStatus string, paymentStatus string) error {
	return nil
}
