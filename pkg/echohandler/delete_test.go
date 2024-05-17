package echohandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/JHTohru/eulabschallenge/pkg/product"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type deleterStub struct {
	prd *product.Product
	err error
}

func (d *deleterStub) Delete(_ context.Context, _ uuid.UUID) (*product.Product, error) {
	return d.prd, d.err
}

func TestDeleteHandler(t *testing.T) {
	errFake := fmt.Errorf("fake error")
	id := uuid.New()
	now := time.Now()
	prd := &product.Product{
		ID:          id,
		Name:        "lorem ipsum",
		Description: "dolor sit amet",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	prdMarsh, err := json.Marshal(prd)
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		id         string
		deletePrd  *product.Product
		deleteErr  error
		statusWant int
		bodyWant   string
		errWant    error
	}{
		"it must err for a malformed id path parameter": {
			id:         "",
			deletePrd:  nil,
			deleteErr:  nil,
			statusWant: http.StatusOK,
			bodyWant:   "",
			errWant:    errMalformedIDPathParam,
		},
		"it must err for a failing delete": {
			id:         id.String(),
			deletePrd:  nil,
			deleteErr:  errFake,
			statusWant: http.StatusOK,
			bodyWant:   "",
			errWant:    errFake,
		},
		"happy path": {
			id:         id.String(),
			deletePrd:  prd,
			deleteErr:  nil,
			statusWant: http.StatusOK,
			bodyWant:   string(prdMarsh),
			errWant:    nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			deleter := &deleterStub{test.deletePrd, test.deleteErr}
			dh := &DeleteHandler{deleter}
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)
			c.SetParamNames("id")
			c.SetParamValues(test.id)

			err := dh.Handle(c)

			respBody := resp.Body.String()
			respBody = strings.TrimSuffix(respBody, "\n")
			if respBody != test.bodyWant {
				t.Errorf("unexpected response body; got %q, want %q", respBody, test.bodyWant)
			}

			if resp.Code != test.statusWant {
				t.Errorf("unexpected response status code; got %d, want %d", resp.Code, test.statusWant)
			}

			if !errors.Is(err, test.errWant) {
				t.Errorf("unexpected error; got %v, want %v", err, test.errWant)
			}
		})
	}
}
