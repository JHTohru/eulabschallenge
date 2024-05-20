package echohandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type Lister interface {
	List(ctx context.Context, pageSize, pageNumber int) (prds []*product.Product, pagesTotal int, err error)
}

type productsPage struct {
	Products   []*product.Product `json:"products"`
	PagesTotal int                `json:"pages_total"`
}

type ListHandler struct {
	lister Lister
}

func NewListHandler(l Lister) *ListHandler {
	return &ListHandler{
		lister: l,
	}
}

func (lh *ListHandler) Handle(c echo.Context) error {
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		return errMalformedPageSizeQueryParam
	}

	pageNumber, err := strconv.Atoi(c.QueryParam("page_number"))
	if err != nil {
		return errMalformedPageNumberQueryParam
	}

	ctx := c.Request().Context()
	prds, pagesTotal, err := lh.lister.List(ctx, pageSize, pageNumber)
	if err != nil {
		return err
	}
	if prds == nil {
		prds = make([]*product.Product, 0)
	}

	page := productsPage{prds, pagesTotal}
	return c.JSON(http.StatusOK, page)
}
