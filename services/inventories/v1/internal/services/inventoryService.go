package services

import (
	"context"
	"encoding/json"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/repository"
	"inventories/v1/proto/Inventory"
	"log"
	"package/rabbitmq/publisher"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type inventoryServer struct {
	tracer        trace.Tracer
	inventoryRepo repository.InventoryRepository
	publisher     publisher.EventPublisher
	Inventory.UnimplementedInventoryServiceServer
}

func NewInventoryServer(inventoryRepo repository.InventoryRepository, publisher publisher.EventPublisher, tracer trace.Tracer) Inventory.InventoryServiceServer {
	return &inventoryServer{
		inventoryRepo: inventoryRepo,
		publisher:     publisher,
	}
}

func (s *inventoryServer) AddInventory(ctx context.Context, in *Inventory.AddInvenRequest) (*Inventory.AddInvenResponse, error) {

	tracer := otel.Tracer("inventory-service")
	addCtx, addSpan := tracer.Start(ctx, "AddInventory")
	productID, err := s.inventoryRepo.AddInventory(addCtx, &constant.Inventory{
		StoreID:        in.Inventory.StoreID,
		AddBy:          in.Inventory.AddBy,
		Name:           in.Inventory.Name,
		Description:    in.Inventory.Description,
		Price:          in.Inventory.Price,
		AvailableStock: in.Inventory.Stock,
		CategoryID:     in.Inventory.CategoryID,
		ImageURL:       in.Inventory.ImageURL,
		CreatedAt:      time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	addSpan.End()

	rbCtx, rbSpan := tracer.Start(ctx, "AddInventory to OrderProducts")
	payload := map[string]interface{}{
		"store_id":     in.Inventory.StoreID,
		"product_id":   productID,
		"product_name": in.Inventory.Name,
		"price":        in.Inventory.Price,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if err := s.publisher.Publish(rbCtx, body, "inventory.exchange", "inventory.created"); err != nil {
		return nil, err
	}

	log.Println("inventory created publish")
	rbSpan.End()

	response := &Inventory.AddInvenResponse{
		Status: "Inventory added successfully",
	}

	return response, nil
}

func (s *inventoryServer) UpdateInventory(ctx context.Context, in *Inventory.UpdateInvenRequest) (*Inventory.UpdateInvenResponse, error) {
	addCtx, addSpan := s.tracer.Start(ctx, "UpdateInventory")
	if err := s.inventoryRepo.UpdateInventory(addCtx, &constant.Inventory{
		ID:             *in.Inventory.ID,
		StoreID:        in.Inventory.StoreID,
		Name:           in.Inventory.Name,
		Description:    in.Inventory.Description,
		Price:          in.Inventory.Price,
		AvailableStock: in.Inventory.Stock,
		CategoryID:     in.Inventory.CategoryID,
		ImageURL:       in.Inventory.ImageURL,
		UpdatedAt:      time.Now().UTC(),
	}); err != nil {
		return nil, err
	}
	addSpan.End()

	upCtx, upSpan := s.tracer.Start(ctx, "UpdateInventory to OrderProducts")
	payload := map[string]interface{}{
		"store_id":     in.Inventory.StoreID,
		"product_id":   *in.Inventory.ID,
		"product_name": in.Inventory.Name,
		"price":        in.Inventory.Price,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if err := s.publisher.Publish(upCtx, body, "inventory.exchange", "inventory.updated"); err != nil {
		return nil, err
	}
	upSpan.End()

	response := &Inventory.UpdateInvenResponse{
		Status: "Inventory updated successfully",
	}

	return response, nil
}

func (s *inventoryServer) RemoveInventory(ctx context.Context, in *Inventory.RemoveInvenRequest) (*Inventory.RemoveInvenResponse, error) {
	if err := s.inventoryRepo.RemoveInventory(in.UserID, in.StoreID, in.InvetoriesID); err != nil {
		return nil, err
	}

	response := &Inventory.RemoveInvenResponse{
		Status: "Inventory removed successfully",
	}

	return response, nil
}

func (s *inventoryServer) GetInventory(ctx context.Context, in *Inventory.GetInvetoryRequest) (*Inventory.GetInvetoryResponse, error) {
	return nil, nil
}

func (s *inventoryServer) ListInventories(ctx context.Context, in *Inventory.ListInvetoriesRequest) (*Inventory.ListInvetoriesResponse, error) {
	req := &constant.ListInventoryReq{
		StoreID:    in.Fields.StoreID,
		Query:      in.Fields.Query,
		CategoryID: in.Fields.CategoryID,
	}
	pagination := &constant.Pagination{
		Limit:  in.GetPagination().GetLimit(),
		Offset: in.GetPagination().GetOffset(),
	}
	data, err := s.inventoryRepo.ListInventory(req, pagination)
	if err != nil {
		return nil, err
	}

	response := &Inventory.ListInvetoriesResponse{
		Inventory: data,
		Pagination: &Inventory.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
			Total:  &pagination.Total,
		},
	}

	return response, nil
}
