package repository

import (
	"auth-service/v1/internal/constant"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(*constant.User) (*constant.User, error)
	Register(*constant.User) error
}

type authRepository struct {
	gorm *gorm.DB
	sqlx *sqlx.DB
}

func NewAuthRepository(gorm *gorm.DB, sqlx *sqlx.DB) AuthRepository {
	return &authRepository{
		gorm: gorm,
		sqlx: sqlx,
	}
}

func (repo *authRepository) Login(user *constant.User) (*constant.User, error) {
	u := &constant.User{}
	result := repo.gorm.Where("email = ?", user.Email).First(u)
	if result.Error != nil {
		return nil, result.Error
	}

	query := `
		SELECT
			r.id
		FROM users u
		INNER JOIN user_role ur
			ON u.id = ur.user_id
		INNER JOIN role r
			ON ur.role_id = r.id
		WHERE u.id = $1
	`
	var role []int
	args := []interface{}{user.ID}
	err := repo.sqlx.Select(&role, query, args...)
	if err != nil {
		return nil, err
	}

	u.Role = role

	return u, nil
}

func (repo *authRepository) Register(user *constant.User) error {
	result := repo.gorm.Omit("updated_at", "verified").Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
