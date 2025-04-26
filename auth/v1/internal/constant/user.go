package constant

import "time"

type User struct {
	ID          string    `json:"id" db:"id" validate:"required"`
	FirstName   string    `json:"first_name" db:"first_name" validate:"required"`
	LastName    string    `json:"last_name" db:"last_name" validate:"required"`
	Role        []int     `json:"role" db:"role" validate:"required" gorm:"-"`
	Email       string    `json:"email" db:"email" validate:"required,email"`
	PhoneNumber string    `json:"phone_number" db:"phone_number" validate:"required"`
	Verified    bool      `json:"verified" db:"verified" gorm:"column:verified"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Password    string    `json:"password" db:"password" validate:"required"`
}
