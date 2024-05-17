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

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/JHTohru/eulabschallenge/pkg/product"
)

type creatorStub struct {
	prd *product.Product
	err error
}

func (c *creatorStub) Create(_ context.Context, _ *product.Input) (*product.Product, error) {
	return c.prd, c.err
}

func TestCreateHandler(t *testing.T) {
	errFake := fmt.Errorf("fake error")
	now := time.Now()
	prd := &product.Product{
		ID:          uuid.New(),
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
		requestBody string
		createPrd   *product.Product
		createErr   error
		statusWant  int
		bodyWant    string
		errWant     error
	}{
		"it must err for a malformed request body": {
			requestBody: "",
			createPrd:   nil,
			createErr:   nil,
			statusWant:  http.StatusOK,
			bodyWant:    "",
			errWant:     errMalformedRequestBody,
		},
		"it must err for failing create": {
			requestBody: `{"name":"lorem ipsum","description":"dolor sit amet"}`,
			createPrd:   nil,
			createErr:   errFake,
			statusWant:  http.StatusOK,
			bodyWant:    "",
			errWant:     errFake,
		},
		"happy path": {
			requestBody: `{"name":"lorem ipsum","description":"dolor sit amet"}`,
			createPrd:   prd,
			createErr:   nil,
			statusWant:  http.StatusCreated,
			bodyWant:    string(prdMarsh),
			errWant:     nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			creator := &creatorStub{test.createPrd, test.createErr}
			ch := &CreateHandler{creator}
			e := echo.New()
			reqBody := strings.NewReader(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", reqBody)
			req.Header.Add("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)

			err := ch.Handle(c)

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
