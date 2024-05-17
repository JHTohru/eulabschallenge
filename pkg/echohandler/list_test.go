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

type listerStub struct {
	prds  []*product.Product
	total int
	err   error
}

func (f *listerStub) List(_ context.Context, _, _ int) ([]*product.Product, int, error) {
	return f.prds, f.total, f.err
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestListHandler(t *testing.T) {
	errFake := fmt.Errorf("fake error")
	prds := []*product.Product{
		{
			ID:          uuid.New(),
			Name:        "lorem ipsum",
			Description: "dolor sit amet",
			CreatedAt:   mustParseTime("2024-03-12 12:44:22"),
			UpdatedAt:   mustParseTime("2024-03-12 12:44:22"),
		},
		{
			ID:          uuid.New(),
			Name:        "consectetur",
			Description: "adipiscing elit",
			CreatedAt:   mustParseTime("2024-04-11 13:21:30"),
			UpdatedAt:   mustParseTime("2024-04-11 14:01:59"),
		},
		{
			ID:          uuid.New(),
			Name:        "donec non",
			Description: "convallis nulla",
			CreatedAt:   mustParseTime("2024-04-15 12:04:51"),
			UpdatedAt:   mustParseTime("2024-05-03 23:11:19"),
		},
	}
	prdsMarsh, err := json.Marshal(prds)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		pageSize       string
		pageNumber     string
		listPrds       []*product.Product
		listTotalPages int
		listErr        error
		statusWant     int
		bodyWant       string
		errWant        error
	}{
		"it must err for a malformed page_size query parameter": {
			pageSize:       "",
			pageNumber:     "1",
			listPrds:       nil,
			listTotalPages: 0,
			listErr:        nil,
			statusWant:     http.StatusOK,
			bodyWant:       "",
			errWant:        errMalformedPageSizeQueryParam,
		},
		"it must err for a malformed page_number query parameter": {
			pageSize:       "3",
			pageNumber:     "",
			listPrds:       nil,
			listTotalPages: 0,
			listErr:        nil,
			statusWant:     http.StatusOK,
			bodyWant:       "",
			errWant:        errMalformedPageNumberQueryParam,
		},
		"it must err for a failing list": {
			pageSize:       "3",
			pageNumber:     "1",
			listPrds:       nil,
			listTotalPages: 0,
			listErr:        errFake,
			statusWant:     http.StatusOK,
			bodyWant:       "",
			errWant:        errFake,
		},
		"happy path": {
			pageSize:       "3",
			pageNumber:     "1",
			listPrds:       prds,
			listTotalPages: 34,
			listErr:        nil,
			statusWant:     http.StatusOK,
			bodyWant:       fmt.Sprintf("{\"products\":%s,\"pages_total\":%d}", prdsMarsh, 34),
			errWant:        nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			lister := &listerStub{test.listPrds, test.listTotalPages, test.listErr}
			lh := &ListHandler{lister}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			q := req.URL.Query()
			q.Add("page_size", test.pageSize)
			q.Add("page_number", test.pageNumber)
			req.URL.RawQuery = q.Encode()
			resp := httptest.NewRecorder()
			c := e.NewContext(req, resp)

			err := lh.Handle(c)

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
