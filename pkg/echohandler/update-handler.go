package echohandler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type Updater interface {
	Update(ctx context.Context, id uuid.UUID, in *product.Input) (*product.Product, error)
}

type UpdateHandler struct {
	updater Updater
}

func NewUpdateHandler(u Updater) *UpdateHandler {
	return &UpdateHandler{
		updater: u,
	}
}

func (uh *UpdateHandler) Handle(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errMalformedIDPathParam
	}

	in := &product.Input{}
	if err := c.Bind(in); err != nil {
		return errMalformedRequestBody
	}

	ctx := c.Request().Context()
	prd, err := uh.updater.Update(ctx, id, in)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, prd)
}
