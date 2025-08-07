package service

import (
	"context"
	"log"
	"package/rabbitmq/publisher"
	"strconv"

	"github.com/stripe/stripe-go/v82/paymentintent"

	"github.com/stripe/stripe-go/v82"
)

type payment struct {
	stripeKey string
	publisher publisher.EventPublisher
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID int, amountPrice float32, userID string) error
}

func NewPaymentService(stripeKey string, publisher publisher.EventPublisher) PaymentService {
	return &payment{
		stripeKey: stripeKey,
		publisher: publisher,
	}
}

func (p *payment) ProcessPayment(ctx context.Context, orderID int, amountPrice float32, userID string) error {
	log.Println("stripeKey: ", p.stripeKey)
	stripe.Key = p.stripeKey

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountPrice * 100)),
		Currency:           stripe.String(string(stripe.CurrencyTHB)),
		PaymentMethod:      stripe.String("pm_card_visa"),
		PaymentMethodTypes: []*string{stripe.String("card")},
		// Confirm:            stripe.Bool(true),
		Metadata: map[string]string{
			"order_id": strconv.Itoa(orderID),
			"user_id":  userID,
		},
	}
	result, err := paymentintent.New(params)
	if err != nil {
		return err
	}

	// godump.Dump(result)
	log.Println("Status: ", result.Status)
	log.Println("ClientSecret: ", result.ClientSecret)
	log.Println("Result payment: ", result)
	return nil
}

func (p *payment) UpdateOrderStatus(order_id int) error {

	return nil
}
