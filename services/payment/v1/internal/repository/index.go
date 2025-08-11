package repository

import (
	"payment/v1/internal/constant"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type PaymentReposiotry interface {
	CreatePayment(*constant.Payment) error
	UpdatePayment(*constant.Payment) error
}

type paymentRepository struct {
	gorm *gorm.DB
	sqlx *sqlx.DB
}

func NewPaymentRepository(gorm *gorm.DB, sqlx *sqlx.DB) PaymentReposiotry {
	return &paymentRepository{
		gorm: gorm,
		sqlx: sqlx,
	}
}
