// build +test
package request

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

// MockResponse to mock a http response whilst unit testing
type MockResponse struct {
	URL          string
	ResponseBody []byte
	StatusCode   int
}

// MockRequest a structure to hold the mocked responses and remember the requests and headers to provide a unit test pinch point
type MockRequest struct {
	CannedResponses []MockResponse
}

// Get makes the HTTP request and writes the results to the io.Writer
func (mr *MockRequest) Get(ctx context.Context, uri string, w io.Writer) error {
	status, response := mr.lookupResponse(uri)

	if status != http.StatusOK {
		return errors.New("Status Not Ok")
	}
	_, err := io.Copy(w, bytes.NewReader(response))
	if err != nil {
		return err
	}
	return nil
}

// AddCannedResponse add a canned response to the list.
func (mr *MockRequest) AddCannedResponse(response MockResponse) {
	mr.CannedResponses = append(mr.CannedResponses, response)
}

func (mr *MockRequest) lookupResponse(url string) (status int, response []byte) {
	for _, r := range mr.CannedResponses {
		if r.URL == url {
			return r.StatusCode, r.ResponseBody
		}
	}
	return http.StatusNotFound, []byte("")
}
