package constant

type Role int32

const (
	Admin Role = iota + 1
	Customer
	Seller
)
