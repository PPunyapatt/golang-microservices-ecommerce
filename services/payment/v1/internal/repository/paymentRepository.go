package repository

import (
	"context"
	"fmt"
	"log"
	"payment/v1/internal/constant"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

func (p *paymentRepository) CreatePayment(ctx context.Context, payment *constant.Payment) error {
	args := []string{"updated_at", "failure_reason", "payment_method"}
	if result := p.gorm.WithContext(ctx).Omit(args...).Create(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *paymentRepository) UpdatePayment(ctx context.Context, payment *constant.Payment, excepts ...string) error {
	if payment.FailureReason == nil {
		excepts = append(excepts, "failure_reason")
	}

	if result := p.gorm.WithContext(ctx).Omit(excepts...).Where("payment_id = ?", payment.PaymentID).Updates(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *paymentRepository) GetPaymentIntentIDbyOrderID(ctx context.Context, orderID int) (string, error) {
	var payment constant.Payment
	result := p.gorm.WithContext(ctx).Where("order_id = ?", orderID).First(&payment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// handle not found
			return "", fmt.Errorf("payment not found for order_id %d", orderID)
		}
	}
	return payment.PaymentID, nil
}

func (p *paymentRepository) CheckPaymentsuccessed(ctx context.Context, orderID int) (bool, error) {
	args := []interface{}{orderID}
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM payments
			where order_id = $1 and status = 'successed'
		)
	`

	var exists bool
	if err := p.sqlx.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		log.Printf("%+v", errors.WithStack(err))
		return false, err
	}
	return exists, nil
}
