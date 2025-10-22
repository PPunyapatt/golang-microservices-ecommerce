package helper

import (
	"errors"

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
	Err          error
}

func (err HttpError) Error() string {
	return err.Message
}

func RespondHttpError(ctx *fiber.Ctx, err error) error {
	var httpError HttpError
	statusCode := fiber.StatusInternalServerError
	ok := errors.As(err, &httpError)
	if ok {
		statusCode = httpError.StatusCode
	}

	if st, ok := status.FromError(httpError.Err); ok {
		return fiber.NewError(statusCode, "rpc - "+st.Message())
	}

	return fiber.NewError(statusCode, err.Error())
}

func NewHttpError(statusCode int, err error) HttpError {
	return HttpError{
		StatusCode: statusCode,
		Message:    err.Error(),
		Err:        err,
	}
}
