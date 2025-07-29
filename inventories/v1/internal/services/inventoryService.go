package services

import (
	"context"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/repository"
	"inventories/v1/proto/Inventory"
)

type inventoryServer struct {
	inventoryRepo repository.InventoryRepository
	Inventory.UnimplementedInventoryServiceServer
}

func NewInventoryServer(inventoryRepo repository.InventoryRepository) Inventory.InventoryServiceServer {
	return &inventoryServer{inventoryRepo: inventoryRepo}
}

func (s *inventoryServer) AddInventory(ctx context.Context, in *Inventory.AddInvenRequest) (*Inventory.AddInvenResponse, error) {
	return nil, nil
}

func (s *inventoryServer) UpdateInventory(ctx context.Context, in *Inventory.UpdateInvenRequest) (*Inventory.UpdateInvenResponse, error) {
	return nil, nil
}

func (s *inventoryServer) RemoveInventory(ctx context.Context, in *Inventory.RemoveInvenRequest) (*Inventory.RemoveInvenResponse, error) {
	return nil, nil
}

func (s *inventoryServer) GetInventory(ctx context.Context, in *Inventory.GetInvetoryRequest) (*Inventory.GetInvetoryResponse, error) {
	return nil, nil
}

func (s *inventoryServer) ListInventory(ctx context.Context, in *Inventory.ListInvetoriesRequest) (*Inventory.ListInvetoriesResponse, error) {
	return nil, nil
}

func (s *inventoryServer) AddCategory(ctx context.Context, in *Inventory.AddCategoryRequest) (*Inventory.AddCategoryResponse, error) {
	catagory := &constant.Category{
		Name: in.GetName(),
	}
	err := s.inventoryRepo.AddCategory(catagory)
	if err != nil {
		return nil, err
	}
	response := &Inventory.AddCategoryResponse{
		Status: "Category added successfully",
	}
	return response, nil
}

func (s *inventoryServer) UpdateCategory(ctx context.Context, in *Inventory.UpdateCategoryRequest) (*Inventory.UpdateCategoryResponse, error) {
	catagory := &constant.Category{
		ID:      in.GetCatagoryID(),
		Name:    in.GetName(),
		StoreID: in.GetStoreID(),
	}

	err := s.inventoryRepo.UpdateCategory(catagory)
	if err != nil {
		return nil, err
	}
	response := &Inventory.UpdateCategoryResponse{
		Status: "Category updated successfully",
	}
	return response, nil
}

func (s *inventoryServer) RemoveCategory(ctx context.Context, in *Inventory.RemoveCatgoryRequest) (*Inventory.RemoveCatgoryResponse, error) {
	return nil, nil
}

func (s *inventoryServer) GetCategory(ctx context.Context, in *Inventory.GetCatgoriesRequest) (*Inventory.GetCatgoriesResponse, error) {
	return nil, nil
}
