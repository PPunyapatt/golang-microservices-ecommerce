package service

import (
	"cart/v1/internal/constant"
	"cart/v1/internal/repository"
	"cart/v1/proto/cart"
	"context"
	"log"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type cartServer struct {
	tracer   trace.Tracer
	Logger   *slog.Logger
	cartRepo repository.CartRepository
	cart.UnimplementedCartServiceServer
}

type cartService struct {
	Logger   *slog.Logger
	tracer   trace.Tracer
	cartRepo repository.CartRepository
}

type CartService interface {
	DeleteCartFromEvent(context.Context, string) error
}

func NewCartServer(cartRepo repository.CartRepository, tracer trace.Tracer, logger *slog.Logger) (CartService, cart.CartServiceServer) {
	return &cartService{
			Logger:   logger,
			cartRepo: cartRepo,
			tracer:   tracer,
		}, &cartServer{
			Logger:   logger,
			cartRepo: cartRepo,
			tracer:   tracer,
		}
}

func (s *cartServer) GetCartByUserID(ctx context.Context, in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	getCartCtx, getCartSpan := s.tracer.Start(ctx, "get cart by userID")
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset,
	}

	// items, err := s.cartRepo.GetItemsByUserID(in.UserId, pagination)
	items, err := s.cartRepo.GetOrCreateCartByUserID(getCartCtx, in.UserId, pagination)
	if err != nil {
		log.Println("Err service: ", err)
		return nil, err
	}

	response := &cart.GetCartResponse{
		StoreItems: items,
		Pagination: &cart.Pagination{
			Limit: pagination.Limit,
			Total: &pagination.TotalCount,
		},
	}

	getCartSpan.End()

	return response, nil
}

func (s *cartServer) AddItemToCart(ctx context.Context, in *cart.AddItemRequest) (*cart.AddToCartResponse, error) {
	// Implementation for adding item to cart

	storeItems := []*constant.StoreItems{}
	for _, store := range in.StoreItems {
		items := []*constant.Item{}
		for _, item := range store.Items {
			items = append(items, &constant.Item{
				ProductID:   int(item.ProductId),
				ProductName: item.ProductName,
				Price:       item.Price,
				Quantity:    int(item.Quantity),
				ImageURL:    item.ImageUrl,
			})
		}
		storeItems = append(storeItems, &constant.StoreItems{
			StoreID: int(store.StoreID),
			Items:   items,
		})
	}

	addItemCtx, addItemSpan := s.tracer.Start(ctx, "add item")

	if err := s.cartRepo.AddItem(addItemCtx, in.UserId, storeItems); err != nil {
		return nil, err
	}

	addItemSpan.End()

	return &cart.AddToCartResponse{
		Status: "Items added to cart successfully",
	}, nil
}

func (s *cartServer) RemoveItem(ctx context.Context, in *cart.RemoveFromCartRequest) (*cart.RemoveFromCartResponse, error) {
	deleteItemCtx, deleteItemSpan := s.tracer.Start(ctx, "Remove item")
	if err := s.cartRepo.RemoveItem(deleteItemCtx, in.UserId, int(in.ItemId)); err != nil {
		return nil, err
	}
	deleteItemSpan.End()

	return &cart.RemoveFromCartResponse{
		Status: "Item removed from cart successfully",
	}, nil
}

func (c *cartService) DeleteCartFromEvent(ctx context.Context, userID string) error {
	deleteCartCtx, deleteCartSpan := c.tracer.Start(ctx, "Remove Cart")
	if err := c.cartRepo.RemoveCart(deleteCartCtx, userID); err != nil {
		return err
	}
	deleteCartSpan.End()

	return nil
}
