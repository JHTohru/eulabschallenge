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

type updaterStub struct {
	prd *product.Product
	err error
}

func (u *updaterStub) Update(_ context.Context, _ uuid.UUID, _ *product.Input) (*product.Product, error) {
	return u.prd, u.err
}

func TestUpdateHandler(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	prd := &product.Product{
		ID:          id,
		Name:        "consectetur",
		Description: "adipiscing elit",
		CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
		UpdatedAt:   now,
	}
	prdMarsh, err := json.Marshal(prd)
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		id          string
		requestBody string
		updatePrd   *product.Product
		updateErr   error
		statusWant  int
		bodyWant    string
		errWant     error
	}{
		"it must err for a malformed id path parameter": {
			id:          "",
			requestBody: `{"name":"consectetur","description":"adipiscing elit"}`,
			updatePrd:   nil,
			updateErr:   nil,
			statusWant:  http.StatusOK,
			bodyWant:    "",
			errWant:     errMalformedIDPathParam,
		},
		"it must err for a malformed request body": {
			id:          id.String(),
			requestBody: "",
			updatePrd:   nil,
			updateErr:   nil,
			statusWant:  http.StatusOK,
			bodyWant:    "",
			errWant:     errMalformedRequestBody,
		},
		"it must err for a failing update": {
			id:          id.String(),
			requestBody: `{"name":"consectetur","description":"adipiscing elit"}`,
			updatePrd:   nil,
			updateErr:   errFake,
			statusWant:  http.StatusOK,
			bodyWant:    "",
			errWant:     errFake,
		},
		"happy path": {
			id:          id.String(),
			requestBody: `{"name":"consectetur","description":"adipiscing elit"}`,
			updatePrd:   prd,
			updateErr:   nil,
			statusWant:  http.StatusOK,
			bodyWant:    string(prdMarsh),
			errWant:     nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			updater := &updaterStub{test.updatePrd, test.updateErr}
			uh := &UpdateHandler{updater}
			e := echo.New()
			reqBody := strings.NewReader(test.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/", reqBody)
			req.Header.Add("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)
			c.SetParamNames("id")
			c.SetParamValues(test.id)

			err := uh.Handle(c)

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
