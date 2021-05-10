package request

import (
	"context"
	"io"
	"net/http"
)

// Requester interface for getting a page
type Requester interface {
	Get(ctx context.Context, uri string, w io.Writer) error
}

// HTTPRequest performs the Get request to download the webpage html
type HTTPRequest struct{}

// Get makes the HTTP request and writes the results to the io.Writer
func (HTTPRequest) Get(ctx context.Context, uri string, w io.Writer) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return err
	}
	c := &http.Client{}
	r, err := c.Do(rq)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = io.Copy(w, r.Body)
	if err != nil {
		return err
	}
	return nil
}
