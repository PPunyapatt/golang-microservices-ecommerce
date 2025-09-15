package repository

import (
	"context"
	"payment/v1/internal/constant"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type PaymentReposiotry interface {
	CreatePayment(context.Context, *constant.Payment) error
	UpdatePayment(context.Context, *constant.Payment, ...string) error
	GetPaymentIntentIDbyOrderID(context.Context, int) (string, error)
	CheckPaymentsuccessed(context.Context, int) (bool, error)
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
