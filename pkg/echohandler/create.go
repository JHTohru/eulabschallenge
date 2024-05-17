package echohandler

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type Creator interface {
	Create(ctx context.Context, in *product.Input) (*product.Product, error)
}

type CreateHandler struct {
	creator Creator
}

func NewCreateHandler(c Creator) *CreateHandler {
	return &CreateHandler{
		creator: c,
	}
}

func (ch *CreateHandler) Handle(c echo.Context) error {
	in := &product.Input{}
	if err := c.Bind(in); err != nil {
		return errMalformedRequestBody
	}
	ctx := c.Request().Context()
	prd, err := ch.creator.Create(ctx, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, prd)
}
