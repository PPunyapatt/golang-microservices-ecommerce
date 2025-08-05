package services

import (
	"context"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/repository"
	"inventories/v1/proto/Inventory"
	"log"
	"time"

	"github.com/goforj/godump"
)

type inventoryServer struct {
	inventoryRepo repository.InventoryRepository
	Inventory.UnimplementedInventoryServiceServer
}

func NewInventoryServer(inventoryRepo repository.InventoryRepository) Inventory.InventoryServiceServer {
	return &inventoryServer{inventoryRepo: inventoryRepo}
}

func (s *inventoryServer) AddInventory(ctx context.Context, in *Inventory.AddInvenRequest) (*Inventory.AddInvenResponse, error) {

	log.Println("User: ", in.Inventory.AddBy)

	if err := s.inventoryRepo.AddInventory(&constant.Inventory{
		StoreID:     in.Inventory.StoreID,
		AddBy:       in.Inventory.AddBy,
		Name:        in.Inventory.Name,
		Description: in.Inventory.Description,
		Price:       in.Inventory.Price,
		Stock:       in.Inventory.Stock,
		CategoryID:  in.Inventory.CategoryID,
		ImageURL:    in.Inventory.ImageURL,
		CreatedAt:   time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

	response := &Inventory.AddInvenResponse{
		Status: "Inventory added successfully",
	}

	return response, nil
}

func (s *inventoryServer) UpdateInventory(ctx context.Context, in *Inventory.UpdateInvenRequest) (*Inventory.UpdateInvenResponse, error) {
	output := godump.DumpStr(in.Inventory)
	log.Println("str", output)
	if err := s.inventoryRepo.UpdateInventory(&constant.Inventory{
		ID:          *in.Inventory.ID,
		StoreID:     in.Inventory.StoreID,
		Name:        in.Inventory.Name,
		Description: in.Inventory.Description,
		Price:       in.Inventory.Price,
		Stock:       in.Inventory.Stock,
		CategoryID:  in.Inventory.CategoryID,
		ImageURL:    in.Inventory.ImageURL,
		UpdatedAt:   time.Now().UTC(),
	}); err != nil {
		return nil, err
	}

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
