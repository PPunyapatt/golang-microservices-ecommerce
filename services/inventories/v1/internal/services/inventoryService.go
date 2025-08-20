package services

import (
	"context"
	"encoding/json"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/repository"
	"inventories/v1/proto/Inventory"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type inventoryServer struct {
	tracer        trace.Tracer
	inventoryRepo repository.InventoryRepository
	publisher     publisher.EventPublisher
	Inventory.UnimplementedInventoryServiceServer
}

type inventoryService struct {
	tracer        trace.Tracer
	inventoryRepo repository.InventoryRepository
	publisher     publisher.EventPublisher
}

type InventoryServie interface {
	ReserveStock(context.Context, *constant.Order) error
	CutStock(context.Context, []*constant.Item) error
	ReleaseStock(context.Context, []*constant.Item) error
}

func NewInventoryServer(inventoryRepo repository.InventoryRepository, publisher publisher.EventPublisher, tracer trace.Tracer) (Inventory.InventoryServiceServer, InventoryServie) {
	return &inventoryServer{
			inventoryRepo: inventoryRepo,
			publisher:     publisher,
			tracer:        tracer,
		}, &inventoryService{
			inventoryRepo: inventoryRepo,
			publisher:     publisher,
			tracer:        tracer,
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

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(rbCtx, rabbitmq.AMQPHeaderCarrier(headers))
	if err := s.publisher.Publish(
		rbCtx,
		body,
		"inventory.exchange",
		"inventory.created",
		headers,
	); err != nil {
		return nil, err
	}

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

	response := &Inventory.UpdateInvenResponse{
		Status: "Inventory updated successfully",
	}

	if in.Inventory.StoreID == nil && in.Inventory.ID == nil && in.Inventory.Name == nil && in.Inventory.Price == nil {
		return response, nil
	}

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

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(ctx, rabbitmq.AMQPHeaderCarrier(headers))

	if err := s.publisher.Publish(
		upCtx,
		body,
		"inventory.exchange",
		"inventory.updated",
		headers,
	); err != nil {
		return nil, err
	}
	upSpan.End()

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

// -------------------------------- Service --------------------------------

func (s *inventoryService) ReserveStock(ctx context.Context, order *constant.Order) error {
	err := s.inventoryRepo.ReserveStock(ctx, order.Items)

	routingKey := "inventory.reserved"
	if err != nil {
		if err.Error() == "The product is out of stock." {
			routingKey = "inventory.failed"
		} else {
			return err
		}
	}

	payload := map[string]interface{}{
		"order_id":    order.OrderID,
		"total_price": order.TotalPrice,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(ctx, rabbitmq.AMQPHeaderCarrier(headers))

	if err = s.publisher.Publish(
		ctx,
		body,
		"inventory.exchange",
		routingKey,
		headers,
	); err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) CutStock(ctx context.Context, items []*constant.Item) error {
	if err := s.inventoryRepo.CutStock(ctx, items); err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) ReleaseStock(ctx context.Context, items []*constant.Item) error {
	if err := s.inventoryRepo.ReleaseStock(ctx, items); err != nil {
		return err
	}

	return nil
}
