package helper

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/status"
)

const (
	ErrorCodeUniqueViolation     = "23505"
	ErrorCodeForeignKeyViolation = "23503"
	ErrorCodeNotNullViolation    = "23502"

	ErrorMessageConflict       = "Conflict"
	ErrorMessageBadRequest     = "Bad request"
	ErrorMessageUnauthorized   = "Unauthorized"
	ErrorMessageNotfound       = "Not found"
	ErrorMessageForbidden      = "Forbidden"
	ErrorMessageInternalServer = "Internal server error"
)

type HttpError struct {
	StatusCode   int
	Message      string
	ErrorMessage string
}

func RespondHttpError(ctx *fiber.Ctx, err error) error {
	if st, ok := status.FromError(err); ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": st.Message(),
		})
	}

	// Respond with a generic error message
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func NewHttpError(statusCode int, message *string) HttpError {
	if message == nil {
		defaultMessage := getDefaultStatusMessage(statusCode)
		message = &defaultMessage
	}
	return HttpError{
		StatusCode: statusCode,
		Message:    *message,
	}
}

func getDefaultStatusMessage(statusCode int) string {
	message := ""
	switch statusCode {
	case http.StatusConflict:
		message = ErrorMessageConflict
	case http.StatusBadRequest:
		message = ErrorMessageBadRequest
	case http.StatusUnauthorized:
		message = ErrorMessageUnauthorized
	case http.StatusNotFound:
		message = ErrorMessageNotfound
	case http.StatusForbidden:
		message = ErrorMessageForbidden
	default:
		message = ErrorMessageInternalServer
	}

	return message
}
