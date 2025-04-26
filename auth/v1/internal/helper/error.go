package helper

import (
	"auth-service/v1/internal/constant"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

type HttpError struct {
	StatusCode   int
	ErrorMessage string
}

func (err HttpError) Error() string {
	return err.ErrorMessage
}

func ResponseHttpError(ctx *fiber.Ctx, err error) error {
	var httpErr HttpError
	ok := errors.As(err, &httpErr)
	log.Println("ok: ", ok)
	if ok {
		return ctx.Status(httpErr.StatusCode).JSON(constant.StatusResponse{
			StatusCode: httpErr.StatusCode,
			Message:    httpErr.ErrorMessage,
		})
	}

	return errors.New("internal server error")
}

func NewHttpErrorWithDetail(code int, err error) HttpError {
	return HttpError{
		StatusCode:   code,
		ErrorMessage: err.Error(),
	}
}
