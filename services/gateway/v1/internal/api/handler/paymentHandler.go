package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/v1/internal/constant"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/payment"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentmethod"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (c *ApiHandler) StripeWebhook(ctx *fiber.Ctx) error {
	body := ctx.Body()
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := ctx.Get("Stripe-Signature")

	// Verify event
	event, err := webhook.ConstructEvent(body, sigHeader, endpointSecret)
	if err != nil {
		return helper.RespondHttpError(ctx, fmt.Errorf("⚠️ Webhook signature verification failed: %v", err))
	}

	var pi stripe.PaymentIntent
	if err = json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Println("Error parsing payment intent: ", err)
		return helper.RespondHttpError(ctx, err)
	}
	metadata := pi.Metadata

	var paymentType string
	if pi.PaymentMethod != nil {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		paymentMethodID := pi.PaymentMethod.ID
		pm, err := paymentmethod.Get(paymentMethodID, nil)
		if err != nil {
			log.Println("Error fetching payment method details:", err)
			return helper.RespondHttpError(ctx, err)

		}
		paymentType = *stripe.String(pm.Type)
	} else {
		log.Println("No PaymentMethod attached to PaymentIntent")
	}

	paymentReq := &payment.StripeWebhookRequest{
		PaymentIntentID: pi.ID,
		EventType:       string(event.Type),
		Currency:        string(pi.Currency),
		MethodType:      paymentType,
		Metadata:        metadata,
	}
	if pi.LastPaymentError != nil {
		log.Printf("Payment failed: %s - %s\n", pi.LastPaymentError.Code, pi.LastPaymentError.Err)
		paymentErr := pi.LastPaymentError.Err.Error()
		paymentReq.ErrorMessage = &paymentErr
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	_, err = c.PaymentSvc.StripeWebhook(context_, paymentReq)
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
	})
}

func (c *ApiHandler) Paid(ctx *fiber.Ctx) error {
	request, err := helper.ParseAndValidateRequest(ctx, &constant.PaymentRequest{})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	userID, ok := ctx.Locals("UserID").(string)
	if !ok {
		return helper.RespondHttpError(ctx, errors.New("user ID not found in context"))
	}

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	_, err = c.PaymentSvc.Paid(context_, &payment.PaymentRequest{
		Amount:   request.AmountPrice,
		Currency: request.Currency,
		OrderID:  int32(request.OrderID),
		UserID:   userID,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
	})
}
