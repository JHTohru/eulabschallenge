package echohandler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

func HTTPErrorHandler(err error, c echo.Context) {
	switch {
	case errors.As(err, new(*malformationError)):
		c.String(http.StatusBadRequest, err.Error())
	case errors.As(err, new(*product.ValidationError)):
		c.String(http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, product.ErrNotFound):
		c.String(http.StatusNotFound, err.Error())
	default:
		c.Logger().Error(err)
		c.String(http.StatusInternalServerError, "internal error")
	}
}
