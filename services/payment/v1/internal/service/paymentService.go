package service

import (
	"context"
	"log"
	"package/rabbitmq/publisher"
	"payment/v1/proto/payment"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v82/paymentintent"
	"go.opentelemetry.io/otel"

	"github.com/stripe/stripe-go/v82"
)

type paymentService struct {
	stripeKey string
	publisher publisher.EventPublisher
}

type paymentServiceRPC struct {
	payment.UnimplementedPaymentServiceServer
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, orderID int, amountPrice float32, userID string) error
	// StripeWebhook(context.Context, *payment.StripeWebhookRequest) (*payment.Empty, error)
}

func NewPaymentService(stripeKey string, publisher publisher.EventPublisher) (PaymentService, payment.PaymentServiceServer) {
	return &paymentService{
		stripeKey: stripeKey,
		publisher: publisher,
	}, &paymentServiceRPC{}
}

func (p *paymentService) ProcessPayment(ctx context.Context, orderID int, amountPrice float32, userID string) error {
	log.Println("stripeKey: ", p.stripeKey)
	stripe.Key = p.stripeKey

	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(amountPrice * 100)),
		Currency:           stripe.String(string(stripe.CurrencyTHB)),
		PaymentMethod:      stripe.String("pm_card_visa"), // Recieve from frontend
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
	log.Println("PaymentIntentID: ", result.ID)
	log.Println("Status: ", result.Status)
	log.Println("ClientSecret: ", result.ClientSecret)
	// log.Println("Result payment: ", result)
	return nil
}

func (p *paymentServiceRPC) StripeWebhook(ctx context.Context, in *payment.StripeWebhookRequest) (*payment.Empty, error) {
	switch in.EventType {
	case "payment_intent.succeeded":
	case "payment_intent.payment_failed":
	default:
	}

	tracer := otel.Tracer("payment-service")
	_, eventSpan := tracer.Start(ctx, "EventType")
	time.Sleep(1 * time.Second)
	log.Println("Event Type: ", in.EventType)
	log.Println("Metadata: ", in.Metadata)
	eventSpan.End()

	return nil, nil
}

func (p *paymentService) UpdateOrderStatus(order_id int) error {

	return nil
}
