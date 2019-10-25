package server

import (
	"bytes"
	"net/http"
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

}
