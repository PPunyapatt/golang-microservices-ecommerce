package constant

type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Pasword string `json:"password"`
	Roles   []int32
}

type UserRegister struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Store struct {
	ID    int32  `gorm:"column:id"`
	Name  string `gorm:"column:name"`
	Owner string `gorm:"column:owner"`
}
