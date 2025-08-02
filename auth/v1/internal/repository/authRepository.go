package repository

import (
	"auth-service/v1/internal/constant"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Login(context.Context, *constant.User) (*constant.User, error)
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

func (repo *authRepository) Login(ctx context.Context, user *constant.User) (*constant.User, error) {
	// u := &constant.User{}
	// result := repo.gorm.WithContext(ctx).Where("email = ?", user.Email).First(u)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	// log.Println("u.id: ", u.ID)

	// query := `
	// 	SELECT
	// 		r.id
	// 	FROM users u
	// 	INNER JOIN user_role ur
	// 		ON u.id = ur.user_id
	// 	INNER JOIN roles r
	// 		ON ur.role_id = r.id
	// 	WHERE u.id = $1
	// `

	query := `
		SELECT
            ARRAY_AGG(r.id) AS role_ids,
            u.id,
            u.password_hash
		FROM users u
		INNER JOIN user_role ur
			ON u.id = ur.user_id
		INNER JOIN roles r
			ON ur.role_id = r.id
		WHERE u.email = $1
        GROUP BY u.id, u.password_hash
	`

	args := []interface{}{user.Email}
	// if err := repo.sqlx.QueryRowxContext(ctx, query, args...).StructScan(user); err != nil {
	// 	return nil, err
	// }
	var u constant.User
	var roles pq.Int64Array
	err := repo.sqlx.QueryRowxContext(ctx, query, args...).Scan(&roles, &u.ID, &u.Password)
	if err != nil {
		log.Println("error: ", err.Error())
		return nil, err
	}

	u.Roles = make([]int32, len(roles))
	for _, role := range roles {
		u.Roles = append(u.Roles, int32(role))
	}

	return &u, nil
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
