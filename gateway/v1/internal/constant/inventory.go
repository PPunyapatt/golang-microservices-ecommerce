package constant

type Inventories struct {
	ID          *int32   `json:"product_id"`
	StoreID     *int32   `json:"store_id"`
	AddBy       *string  `json:"add_by"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int32   `json:"stock"`
	CategoryID  *int32   `json:"category_id"`
	ImageURL    *string  `json:"image_url"`
}

type Category struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	StoreID int32  `json:"store_id"`
}

type GetCategoryReq struct {
	StoreID    *int32  `json:"store_id" query:"store_id"`
	Query      *string `json:"search" query:"query"`
	CategoryID *int32  `json:"category_id" query:"category_id"`
}

type RemoveInventoryReq struct {
	StoreID     int32 `json:"store_id" params:"store_id"`
	InventoryID int32 `json:"inventory_id" params:"inventory_id"`
}
