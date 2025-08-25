package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"payment/v1/internal/constant"
	"payment/v1/internal/repository"
	"payment/v1/proto/payment"
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/refund"
	"go.opentelemetry.io/otel"

	"github.com/stripe/stripe-go/v82"
)

type paymentService struct {
	stripeKey   string
	publisher   publisher.EventPublisher
	paymentRepo repository.PaymentReposiotry
}

type paymentServiceRPC struct {
	stripeKey      string
	paymentService PaymentService
	paymentRepo    repository.PaymentReposiotry
	publisher      publisher.EventPublisher
	payment.UnimplementedPaymentServiceServer
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, payment *constant.PaymentRequest) error
	ProcessRefund(ctx context.Context, paymentIntentID string) error
}

func NewPaymentService(stripeKey string, paymentRepo repository.PaymentReposiotry, publisher publisher.EventPublisher) (PaymentService, payment.PaymentServiceServer) {
	service := &paymentService{
		stripeKey:   stripeKey,
		publisher:   publisher,
		paymentRepo: paymentRepo,
	}
	return service, &paymentServiceRPC{
		paymentRepo:    paymentRepo,
		stripeKey:      stripeKey,
		paymentService: service,
		publisher:      publisher,
	}
}

func (p *paymentServiceRPC) StripeWebhook(ctx context.Context, in *payment.StripeWebhookRequest) (*payment.Empty, error) {
	tracer := otel.Tracer("payment-service")
	_, eventSpan := tracer.Start(ctx, "EventType")
	defer eventSpan.End()

	orderID, err := strconv.Atoi(in.Metadata["order_id"])
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	paymentData := &constant.Payment{
		PaymentID:     in.PaymentIntentID,
		OrderID:       orderID,
		PaymentMethod: in.MethodType,
		UpdatedAt:     time.Now().UTC(),
	}

	log.Println("Payment Event: ", in.EventType)
	switch in.EventType {
	case "payment_intent.succeeded":
		paymentData.Status = "successed"
	case "payment_intent.payment_failed":
		paymentData.Status = "failed"
		paymentData.FailureReason = in.ErrorMessage
	default:
		return nil, nil
	}

	if err := p.paymentRepo.UpdatePayment(paymentData); err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"order_id":     orderID,
		"user_id":      in.Metadata["user_id"],
		"order_source": in.Metadata["order_source"],
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	headers := amqp091.Table{}
	otel.GetTextMapPropagator().Inject(ctx, rabbitmq.AMQPHeaderCarrier(headers))

	routingKey := "payment." + paymentData.Status
	if err = p.publisher.Publish(
		ctx,
		body,
		"payment.exchange",
		routingKey,
		headers,
	); err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *paymentService) ProcessPayment(ctx context.Context, payment *constant.PaymentRequest) error {
	stripe.Key = p.stripeKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.TotalPrice * 100)),
		Currency: stripe.String("thb"),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
			stripe.String("promptpay"),
		},
		Metadata: map[string]string{
			"order_id":     strconv.Itoa(int(payment.OrderID)),
			"user_id":      payment.UserID,
			"order_source": payment.OrderSource,
		},
	}

	result, err := paymentintent.New(params)
	if err != nil {
		return err
	}

	log.Println("payment ID: ", result.ID)

	paymentData := &constant.Payment{
		PaymentID: result.ID,
		Status:    "pending",
		Amount:    payment.TotalPrice,
		OrderID:   int(payment.OrderID),
		Currency:  string(stripe.CurrencyTHB),
		CreatedAt: time.Now().UTC(),
	}

	if err := p.paymentRepo.CreatePayment(paymentData); err != nil {
		return err
	}

	return nil
}

func (p *paymentService) ProcessRefund(ctx context.Context, paymentIntentID string) error {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return err
	}

	// ใน SDK ใหม่ ใช้ LatestCharge แทน Charges.Data[0]
	if pi.LatestCharge == nil {
		return fmt.Errorf("no charge found for payment intent %s", paymentIntentID)
	}

	chargeID := pi.LatestCharge.ID

	params := &stripe.RefundParams{Charge: stripe.String(chargeID)}
	result, err := refund.New(params)
	if err != nil {
		log.Fatalf("refund failed: %v", err)
	}

	fmt.Printf("Refund created: %s, Status: %s\n", result.ID, result.Status)
	return nil
}
