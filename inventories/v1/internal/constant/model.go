package constant

type Inventory struct {
	ID          int32    `json:"id" db:"id"` // ID ยังใช้ int32 ปกติ
	StoreID     *int32   `json:"store_id" db:"store_id"`
	AddBy       *string  `json:"add_by" db:"add_by"`
	Name        *string  `json:"name" db:"name"`
	Description *string  `json:"description" db:"description"`
	Price       *float64 `json:"price" db:"price"`
	Stock       *int32   `json:"stock" db:"stock"` // ✅ รองรับ 0 ได้
	CategoryID  *int32   `json:"category_id" db:"category_id"`
	ImageURL    *string  `json:"image_url" db:"image_url"`
}

type Category struct {
	ID      int32  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	StoreID int32  `json:"store_id" db:"store_id"`
}
