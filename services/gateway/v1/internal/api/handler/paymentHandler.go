package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/v1/internal/helper"
	"gateway/v1/proto/payment"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (c *ApiHandler) StripeWebhook(ctx *fiber.Ctx) error {
	body := ctx.Body()
	// godump.Dump(body)
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	sigHeader := ctx.Get("Stripe-Signature")

	// Verify event
	event, err := webhook.ConstructEvent(body, sigHeader, endpointSecret)
	if err != nil {
		return helper.RespondHttpError(ctx, fmt.Errorf("⚠️ Webhook signature verification failed: %v", err))
	}
	log.Println("Event Type: ", string(event.Type))

	var pi stripe.PaymentIntent
	if err = json.Unmarshal(event.Data.Raw, &pi); err != nil {
		log.Println("Error parsing payment intent: ", err)
		return helper.RespondHttpError(ctx, err)
	}
	metadata := pi.Metadata

	context_, cancel := context.WithTimeout(ctx.UserContext(), time.Second*10)
	defer cancel()

	_, err = c.PaymentSvc.StripeWebhook(context_, &payment.StripeWebhookRequest{
		EventType: string(event.Type),
		Metadata:  metadata,
	})
	if err != nil {
		return helper.RespondHttpError(ctx, err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
	})
}
