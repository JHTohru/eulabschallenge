package echohandler

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type Getter interface {
	Get(ctx context.Context, id uuid.UUID) (*product.Product, error)
}

type GetHandler struct {
	getter Getter
}

func NewGetHandler(g Getter) *GetHandler {
	return &GetHandler{
		getter: g,
	}
}

func (gh *GetHandler) Handle(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errMalformedIDPathParam
	}

	ctx := c.Request().Context()
	prd, err := gh.getter.Get(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, prd)
}
