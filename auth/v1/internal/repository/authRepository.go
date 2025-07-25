package repository

import (
	"auth-service/v1/internal/constant"
	"log"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(*constant.User) (*constant.User, error)
	Register(*constant.User) error
	CreateStore(*constant.Store) error
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

	log.Println("u.id: ", u.ID)

	query := `
		SELECT
			r.id
		FROM users u
		INNER JOIN user_role ur
			ON u.id = ur.user_id
		INNER JOIN roles r
			ON ur.role_id = r.id
		WHERE u.id = $1
	`
	var role []int32
	args := []interface{}{u.ID}
	err := repo.sqlx.Select(&role, query, args...)
	if err != nil {
		return nil, err
	}

	u.Roles = role

	return u, nil
}

func (repo *authRepository) Register(user *constant.User) error {
	result := repo.gorm.Omit("updated_at", "verified").Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *authRepository) CreateStore(store *constant.Store) error {
	result := repo.gorm.Select("name", "owner").Create(&store)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
