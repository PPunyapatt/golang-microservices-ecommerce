package services

import context "context"

type CartServer struct {
}

func NewCartServer() CartServiceServer {
	return CartServer
}

func (s CartServer) AddItem(ctx context.Context, req *AddItemRequest) (*AddItemResponse, error) {

}

// func (s CartServer) RemoveItem(ctx context.Context, req *RemoveItemRequest) (*RemoveItemResponse, error) {

// }

// func (s CartServer) GetCart(ctx context.Context, req *GetCartRequest) (*GetCartResponse, error) {

// }

func (s CartServer) mustEmbedUnimplementedCartServiceServer() {}
