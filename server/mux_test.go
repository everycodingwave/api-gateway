package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/everycodingwave/api-gateway/cache"
)

var (
	fakeGetRespBody  = []byte("{\"fake_key1\": \"fake_properties_1\"}")
	fakePostRespBody = []byte("{\"contact_id\": \"fake_contact_id_123\"}")
)

func dummyProxyHTTP(method string, url string, errmsg string, w http.ResponseWriter, r *http.Request, handler respHandlerFunc) {
	var callbackData []byte

	if method == "GET" {
		callbackData = fakeGetRespBody
	}

	if method == "POST" || method == "PUT" {
		callbackData = fakePostRespBody
	}

	handler(callbackData)

	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, bytes.NewReader(callbackData)); err != nil {
		log.Printf("%s, copy resp failed: %+v", errmsg, err)
	}
}

type dummyCache struct {
	db map[string]string
}

func newDummyCache() cache.Cache {
	return &dummyCache{
		db: make(map[string]string),
	}
}

func (c *dummyCache) Set(key string, value string, expire time.Duration) error {
	c.db[key] = value
	return nil
}

func (c *dummyCache) Get(key string) (string, error) {
	val, ok := c.db[key]
	if !ok {
		return "", cache.KeyNotExisted
	}

	return val, nil
}

func (c *dummyCache) Del(key string) error {
	delete(c.db, key)
	return nil
}

func TestGetContact(t *testing.T) {
	cac := newDummyCache()
	// server.New actually, same package
	srv := New(cac, dummyProxyHTTP)

	// Testing GET endpoint
	req, err := http.NewRequest("GET", contactAPIURL+"/fake_contact_id_123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// using http response recorder to record down the response and then check
	// have to go thourgh mux cause mux will set context for url parameters like contact_id
	router := srv.getRouter()
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// testing http response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("get contact handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(fakeGetRespBody)
	if rr.Body.String() != expected {
		t.Errorf("get contact handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// testing cache
	val, err := cac.Get("fake_contact_id_123")
	if err != nil || val != string(fakeGetRespBody) {
		t.Errorf("expect cache content fake_contact_id_123=%s, but got %s %v", string(fakeGetRespBody), val, err)
	}
}

func TestCreateContact(t *testing.T) {
	cac := newDummyCache()
	// for further testing cache purging logic
	cac.Set("fake_contact_id_123", "random_data", 0)
	// server.New actually, same package
	srv := New(cac, dummyProxyHTTP)
	// Testing GET endpoint
	req, err := http.NewRequest("POST", contactAPIURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// using http response recorder to record down the response and then check
	// have to go thourgh mux cause mux will set context for url parameters like contact_id
	router := srv.getRouter()
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// testing http response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("create contact handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := string(fakePostRespBody)
	if rr.Body.String() != expected {
		t.Errorf("create contact handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// testing cache
	val, err := cac.Get("fake_contact_id_123")
	if err != cache.KeyNotExisted {
		t.Errorf("expect cache content fake_contact_id_123 being purged, but got %s %v", val, err)
	}
}
