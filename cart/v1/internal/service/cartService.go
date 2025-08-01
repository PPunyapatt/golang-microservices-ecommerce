package service

import (
	"cart/v1/internal/constant"
	"cart/v1/internal/repository"
	"cart/v1/proto/cart"
	"context"
)

type CartService interface {
}

type cartServer struct {
	cartRepo repository.CartRepository
	cart.UnimplementedCartServiceServer
}

func NewCartServer(cartRepo repository.CartRepository) cart.CartServiceServer {
	return &cartServer{cartRepo: cartRepo}
}

func (s *cartServer) GetCartByUserID(ctx context.Context, in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset}

	items, err := s.cartRepo.GetItemsByUserID(in.UserId, pagination)

	if err != nil {
		return nil, err
	}

	response := &cart.GetCartResponse{
		Items: items,
		Pagination: &cart.Pagination{
			Limit: pagination.Limit,
			Total: &pagination.TotalCount,
		},
	}

	return response, nil
}

func (s *cartServer) AddItemToCart(ctx context.Context, in *cart.AddItemRequest) (*cart.AddToCartResponse, error) {
	// Implementation for adding item to cart
	items := []*constant.Item{}
	for _, item := range in.Items {
		items = append(items, &constant.Item{
			ProductID:   int(item.ProductId),
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			Price:       item.Price,
			ImageURL:    item.ImageUrl,
			StoreID:     int(item.StoreID),
		})
	}
	if err := s.cartRepo.AddItem(in.UserId, items); err != nil {
		return nil, err
	}

	return &cart.AddToCartResponse{
		Status: "Items added to cart successfully",
	}, nil
}

func (s *cartServer) RemoveItem(ctx context.Context, in *cart.RemoveFromCartRequest) (*cart.RemoveFromCartResponse, error) {
	if err := s.cartRepo.RemoveItem(in.UserId, int(in.CartId), int(in.CartItemId)); err != nil {
		return nil, err
	}

	return &cart.RemoveFromCartResponse{
		Status: "Item removed from cart successfully",
	}, nil
}
