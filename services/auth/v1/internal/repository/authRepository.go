package repository

import (
	"auth-service/v1/internal/constant"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(context.Context, *constant.User) (*constant.User, error)
	Register(context.Context, *constant.User) error
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

func (repo *authRepository) Login(ctx context.Context, user *constant.User) (*constant.User, error) {
	query := `
		SELECT
			u.id,
            u.password_hash,
			(
				SELECT ARRAY_AGG(r.id)
				FROM user_role ur
				JOIN roles r ON ur.role_id = r.id
				WHERE ur.user_id = u.id) AS role_ids
		FROM users u
		WHERE u.email = $1;
	`

	args := []interface{}{user.Email}
	var u constant.User
	var roles pq.Int64Array
	err := repo.sqlx.QueryRowxContext(ctx, query, args...).Scan(&u.ID, &u.Password, &roles)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return nil, err
	}

	u.Roles = make([]int32, len(roles))
	for i, role := range roles {
		u.Roles[i] = int32(role)
	}

	return &u, nil
}

func (repo *authRepository) Register(ctx context.Context, user *constant.User) error {
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
