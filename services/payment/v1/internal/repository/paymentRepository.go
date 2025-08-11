package repository

import "payment/v1/internal/constant"

func (p *paymentRepository) CreatePayment(payment *constant.Payment) error {
	args := []string{"updated_at", "failure_reason", "payment_method"}
	if result := p.gorm.Omit(args...).Create(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *paymentRepository) UpdatePayment(payment *constant.Payment) error {
	args := []string{"order_id", "amount", "payment_id"}
	if payment.FailureReason == nil {
		args = append(args, "failure_reason")
	}

	if result := p.gorm.Omit(args...).Where("payment_id = ?", payment.PaymentID).Updates(&payment); result.Error != nil {
		return result.Error
	}
	return nil
}
