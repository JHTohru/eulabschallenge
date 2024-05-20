package echohandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type getterStub struct {
	prd *product.Product
	err error
}

func (g *getterStub) Get(_ context.Context, _ uuid.UUID) (*product.Product, error) {
	return g.prd, g.err
}

func TestGetHandler(t *testing.T) {
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
		getPrd     *product.Product
		getErr     error
		statusWant int
		bodyWant   string
		errWant    error
	}{
		"it must err for a malformed id path parameter": {
			id:         "",
			getPrd:     prd,
			getErr:     nil,
			statusWant: http.StatusOK,
			bodyWant:   "",
			errWant:    errMalformedIDPathParam,
		},
		"it must err for a failing get": {
			id:         id.String(),
			getPrd:     nil,
			getErr:     errFake,
			statusWant: http.StatusOK,
			bodyWant:   "",
			errWant:    errFake,
		},
		"happy path": {
			id:         id.String(),
			getPrd:     prd,
			getErr:     nil,
			statusWant: http.StatusOK,
			bodyWant:   string(prdMarsh),
			errWant:    nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			getter := &getterStub{test.getPrd, test.getErr}
			gh := &GetHandler{getter}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)
			c.SetParamNames("id")
			c.SetParamValues(test.id)

			err := gh.Handle(c)

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
