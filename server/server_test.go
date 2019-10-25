package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCopyHTTPHeaders(t *testing.T) {
	src := http.Header{
		"Content-Type":   {"text/html; charset=UTF-8"},
		"Content-Length": {"0"},
	}

	dst := http.Header{}
	copyHeader(dst, src)

	if val := dst.Get("Content-Type"); val != "text/html; charset=UTF-8" {
		t.Errorf("expect Content-Type:text/html; charset=UTF-8, but got %s", val)
	}

	if val := dst.Get("Content-Length"); val != "0" {
		t.Errorf("expect Content-Length 0, but got %s", val)
	}
}

func TestCreateHTTPRequest(t *testing.T) {
	oreq, err := http.NewRequest("GET", "/fake_url", bytes.NewReader([]byte("fake body")))
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}

	oreq.Header.Add(apiAuthHeader, "fake_token")
	nreq, err := createHTTPRequest("GET", "new_fake_url", oreq)

	if err != nil {
		t.Errorf("expect create request without error but got %v", err)
	}

	if nreq.URL.RequestURI() != "new_fake_url" {
		t.Errorf("expect new_fake_url, but got %s", nreq.URL.RequestURI())
	}

	if nreq.Body != oreq.Body {
		t.Errorf("expect request body being copied")
	}

	if val := nreq.Header.Get(apiAuthHeader); val != "fake_token" {
		t.Errorf("expect header %s:fake_token, but got %s", apiAuthHeader, val)
	}

	// testing error case
	oreq.Header.Del(apiAuthHeader)
	nreq, err = createHTTPRequest("GET", "new_fake_url", oreq)
	if err != missingAuthToken {
		t.Errorf("expect throw missing auth token error, but got %v", err)
	}

}

func TestHandleHTTPError(t *testing.T) {
	rr := httptest.NewRecorder()
	handleHTTPError("", errors.New("random error"), rr)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expect http error code %d but got %d", http.StatusInternalServerError, rr.Code)
	}

	rr = httptest.NewRecorder()
	handleHTTPError("", missingAuthToken, rr)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expect http error code %d but got %d", http.StatusBadRequest, rr.Code)
	}
}
