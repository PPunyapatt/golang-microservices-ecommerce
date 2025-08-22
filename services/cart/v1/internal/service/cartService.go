package service

import (
	"cart/v1/internal/constant"
	"cart/v1/internal/repository"
	"cart/v1/proto/cart"
	"context"
	"log"
	"time"
)

type cartServer struct {
	cartRepo repository.CartRepository
	cart.UnimplementedCartServiceServer
}

type cartService struct {
	cartRepo repository.CartRepository
}

type CartService interface {
	DeleteCart(context.Context, string) error
}

func NewCartServer(cartRepo repository.CartRepository) (CartService, cart.CartServiceServer) {
	return &cartService{
			cartRepo: cartRepo,
		}, &cartServer{
			cartRepo: cartRepo,
		}
}

func (s *cartServer) GetCartByUserID(ctx context.Context, in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset,
	}

	log.Println("User Service: ", in.UserId)

	// items, err := s.cartRepo.GetItemsByUserID(in.UserId, pagination)
	items, err := s.cartRepo.GetOrCreateCartByUserID(ctx, in.UserId, pagination)
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
		time_ := time.Now().UTC()
		items = append(items, &constant.Item{
			ProductID:   int(item.ProductId),
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			Price:       item.Price,
			ImageURL:    item.ImageUrl,
			StoreID:     int(item.StoreID),
			CreatedAt:   &time_,
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

func (c *cartService) DeleteCart(ctx context.Context, userID string) error {
	return nil
}
