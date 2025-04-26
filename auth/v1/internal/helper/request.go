package helper

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type ParseOptions struct {
	Params bool
	Query  bool
	Body   bool
}

func ParseAndValidateRequest[T interface{}](context *fiber.Ctx, request *T, options ...ParseOptions) (*T, error) {
	option := ParseOptions{Body: true} // default value
	if len(options) > 0 {
		option = options[0]
	}

	if option.Params {
		if err := context.ParamsParser(request); err != nil {
			log.Printf("params parser error : %+v", errors.WithStack(err))
			return nil, NewHttpErrorWithDetail(http.StatusBadRequest, errors.Wrap(err, "params parse"))
		}
	}

	if option.Query {
		if err := context.QueryParser(request); err != nil {
			log.Printf("query parser error : %+v", errors.WithStack(err))
			return nil, NewHttpErrorWithDetail(http.StatusBadRequest, errors.Wrap(err, "query parse"))
		}
	}

	if option.Body {
		if err := context.BodyParser(request); err != nil {
			log.Printf("body parser error : %+v", errors.WithStack(err))
			return nil, NewHttpErrorWithDetail(http.StatusBadRequest, errors.Wrap(err, "body parse"))
		}
	}

	return request, nil
}
