package echohandler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

func TestHttpErrorHandler(t *testing.T) {
	tests := map[string]struct {
		err        error
		statusWant int
		bodyWant   string
	}{
		"it must respond with the bad request status for a malformation request error": {
			err:        newMalformationError("fake malformation error"),
			statusWant: http.StatusBadRequest,
			bodyWant:   "fake malformation error",
		},
		"it must respond with the unprocessable entity status for validation errors": {
			err:        product.NewValidationError("fake validation error"),
			statusWant: http.StatusUnprocessableEntity,
			bodyWant:   "fake validation error",
		},
		"it must respond with the not found status for validation errors": {
			err:        product.ErrNotFound,
			statusWant: http.StatusNotFound,
			bodyWant:   "product not found",
		},
		"it must respond with the internal server error status for an unexpected errors": {
			err:        errFake,
			statusWant: http.StatusInternalServerError,
			bodyWant:   "internal error",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)

			httpErrorHandler(test.err, c)

			respBody := resp.Body.String()
			respBody = strings.TrimSuffix(respBody, "\n")
			if respBody != test.bodyWant {
				t.Errorf("unexpected response body; got %q, want %q", respBody, test.bodyWant)
			}

			if resp.Code != test.statusWant {
				t.Errorf("unexpected response status code; got %d, want %d", resp.Code, test.statusWant)
			}
		})
	}
}
