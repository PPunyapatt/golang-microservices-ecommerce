package services

import (
	"context"
	"inventories/v1/internal/constant"
	"inventories/v1/internal/repository"
	"inventories/v1/proto/Inventory"
	"log"
	"time"
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
		StoreID:     &in.Inventory.StoreID,
		AddBy:       &in.Inventory.AddBy,
		Name:        &in.Inventory.Name,
		Description: &in.Inventory.Description,
		Price:       &in.Inventory.Price,
		Stock:       &in.Inventory.Stock,
		CategoryID:  &in.Inventory.CategoryID,
		ImageURL:    &in.Inventory.ImageURL,
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
	if err := s.inventoryRepo.UpdateInventory(&constant.Inventory{
		ID:          in.Inventory.ID,
		StoreID:     &in.Inventory.StoreID,
		Name:        &in.Inventory.Name,
		Description: &in.Inventory.Description,
		Price:       &in.Inventory.Price,
		Stock:       &in.Inventory.Stock,
		CategoryID:  &in.Inventory.CategoryID,
		ImageURL:    &in.Inventory.ImageURL,
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
		ID:      in.GetCategoryID(),
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

func (s *inventoryServer) ListCategories(ctx context.Context, in *Inventory.ListCategoriesRequest) (*Inventory.ListCategoriesResponse, error) {
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset,
	}

	log.Println("ListCategories request pagination:", pagination)

	categories, err := s.inventoryRepo.ListCategories(in.GetStoreID(), in.GetSearch(), pagination)

	if err != nil {
		return nil, err
	}

	response := &Inventory.ListCategoriesResponse{
		Catagories: categories,
		Pagination: &Inventory.Pagination{
			Limit:  pagination.Limit,
			Offset: pagination.Offset,
			Total:  &pagination.Total,
		},
	}

	log.Println("ListCategories response pagination:", response.Pagination)
	return response, nil
}
