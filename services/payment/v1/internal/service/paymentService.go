package service

import (
	"context"
	"encoding/json"
	"log"
	"package/rabbitmq"
	"package/rabbitmq/publisher"
	"payment/v1/internal/constant"
	"payment/v1/internal/repository"
	"payment/v1/proto/payment"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/stripe/stripe-go/v82"
)

type paymentService struct {
	stripeKey   string
	publisher   publisher.EventPublisher
	paymentRepo repository.PaymentReposiotry
	tracer      trace.Tracer
}

type paymentServiceRPC struct {
	stripeKey      string
	paymentService PaymentService
	paymentRepo    repository.PaymentReposiotry
	publisher      publisher.EventPublisher
	tracer         trace.Tracer
	payment.UnimplementedPaymentServiceServer
}

type PaymentService interface {
	ProcessPayment(ctx context.Context, payment *constant.PaymentRequest) error
	CancelPayment(ctx context.Context, orderID int) error
}

func NewPaymentService(stripeKey string, paymentRepo repository.PaymentReposiotry, publisher publisher.EventPublisher, tracer trace.Tracer) (PaymentService, payment.PaymentServiceServer) {
	service := &paymentService{
		stripeKey:   stripeKey,
		publisher:   publisher,
		paymentRepo: paymentRepo,
		tracer:      tracer,
	}
	return service, &paymentServiceRPC{
		paymentRepo:    paymentRepo,
		stripeKey:      stripeKey,
		paymentService: service,
		publisher:      publisher,
		tracer:         tracer,
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

	if err := p.paymentRepo.UpdatePayment(context.Background(), paymentData, "order_id", "amount", "payment_id"); err != nil {
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
	paymentCtx, paymentSpan := p.tracer.Start(ctx, "process payment")
	stripe.Key = p.stripeKey

	_, stripeSpan := p.tracer.Start(ctx, "create payment stripe")
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
	stripeSpan.End()

	log.Println("payment ID: ", result.ID)

	paymentData := &constant.Payment{
		PaymentID: result.ID,
		Status:    "pending",
		Amount:    payment.TotalPrice,
		OrderID:   int(payment.OrderID),
		Currency:  string(stripe.CurrencyTHB),
		CreatedAt: time.Now().UTC(),
	}

	if err := p.paymentRepo.CreatePayment(paymentCtx, paymentData); err != nil {
		return err
	}

	paymentSpan.End()

	return nil
}

func (p *paymentService) CancelPayment(ctx context.Context, orderID int) error {
	cancelCtx, cancelSpan := p.tracer.Start(ctx, "cancel payment")
	paymentIntentID, err := p.paymentRepo.GetPaymentIntentIDbyOrderID(ctx, orderID)
	if err != nil {
		return nil
	}

	_, stripeCancelSpan := p.tracer.Start(cancelCtx, "cancel payment")
	params := &stripe.PaymentIntentCancelParams{}
	_, err = paymentintent.Cancel(paymentIntentID, params)
	if err != nil {
		log.Printf("%+v", errors.WithStack(err))
	}
	stripeCancelSpan.End()

	reason := "cancel_payment"
	paymentData := &constant.Payment{
		PaymentID:     paymentIntentID,
		FailureReason: &reason,
		UpdatedAt:     time.Now().UTC(),
		Status:        "cancel",
	}

	exceptUpdates := []string{"payment_id", "order_id", "amount", "currency", "payment_method"}

	err = p.paymentRepo.UpdatePayment(cancelCtx, paymentData, exceptUpdates...)
	if err != nil {
		return err
	}

	cancelSpan.End()

	return nil
}
