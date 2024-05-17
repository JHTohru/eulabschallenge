package echohandler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type Deleter interface {
	Delete(ctx context.Context, id uuid.UUID) (*product.Product, error)
}

type DeleteHandler struct {
	deleter Deleter
}

func NewDeleteHandler(d Deleter) *DeleteHandler {
	return &DeleteHandler{
		deleter: d,
	}
}

func (dh *DeleteHandler) Handle(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errMalformedIDPathParam
	}

	ctx := c.Request().Context()
	prd, err := dh.deleter.Delete(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, prd)
}
