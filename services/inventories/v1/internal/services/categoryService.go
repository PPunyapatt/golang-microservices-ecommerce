package services

import (
	"context"
	"inventories/v1/internal/constant"
	"inventories/v1/proto/Inventory"
	"log"
)

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
	req := &constant.ListInventoryReq{
		StoreID:    in.Fields.StoreID,
		Query:      in.Fields.Query,
		CategoryID: in.Fields.CategoryID,
	}
	pagination := &constant.Pagination{
		Limit:  in.Pagination.Limit,
		Offset: in.Pagination.Offset,
	}

	log.Println("ListCategories request pagination:", pagination)

	categories, err := s.inventoryRepo.ListCategories(req, pagination)

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
