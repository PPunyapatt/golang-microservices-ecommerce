package constant

type Item struct {
	ProductID   int
	PriductName string
	Price       float64
	Quntity     uint32
}

type Cart struct {
	CartID int
	userID string
	Items  []*Item
}
