package service

import (
	"cart/v1/internal/constant"
	"cart/v1/internal/repository"
	"cart/v1/proto/cart"
	"context"
	"log"
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

func (s *cartServer) AddItemToCart(ctx context.Context, in *cart.AddItemRequest) (*cart.AddToCartResponse, error) {
	// Implementation for adding item to cart
	log.Println("AddItemToCart")
	items := []*constant.Item{}
	for _, item := range in.Items {
		items = append(items, &constant.Item{
			ProductID:   int(item.ProductId),
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			Price:       item.Price,
		})
	}
	if err := s.cartRepo.AddItem(in.UserId, items); err != nil {
		return nil, err
	}

	return &cart.AddToCartResponse{
		Status: "Items added to cart successfully",
	}, nil
}
