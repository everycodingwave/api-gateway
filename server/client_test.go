package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpClientTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * httpIOTimeout)
	}))

	defer srv.Close()

	ch := make(chan error)
	go func() {
		_, err := httpClient.Get(srv.URL)
		ch <- err
	}()

	select {
	case <-ch:
	case <-time.After(httpIOTimeout + time.Second):
		t.Errorf("expected http client timeout within %d", httpIOTimeout)
	}
}
