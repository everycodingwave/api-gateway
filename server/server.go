package server

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/everycodingwave/api-gateway/cache"
)

const (
	apiAuthHeader string = "autopilotapikey"
	contactAPIURL string = "https://api2.autopilothq.com/v1/contact"
)

var missingAuthToken = errors.New("missing auth token")

// callback func called after getting backend api response successfully
type respHandlerFunc func(bs []byte)

// the implemention of both endpoint is almost the same, so better extract out the common logic as in another func
// here using a injected callback func is for better testing the apiServer
type httpProxyFunc func(method string, url string, errmsg string, w http.ResponseWriter, r *http.Request, handler respHandlerFunc)

type apiServer struct {
	cac       cache.Cache
	proxyFunc httpProxyFunc
}

func New(c cache.Cache, proxyFunc httpProxyFunc) *apiServer {
	return &apiServer{
		cac:       c,
		proxyFunc: proxyFunc,
	}
}

func (s *apiServer) Start() error {
	log.Printf("api gateway server has started")
	return http.ListenAndServe(":8080", s.getRouter())
}

// ProxyHTTP serves as the core logic of proxying request from gateway server to backend api server.
// It copies gateway request headers and send request to backend api server, and then callback handler func
// inside the callback handler func it does the caching logic
// finally ProxyHTTP will copy backend api response to the client of gateway server.
// so the workflow is:
// client req -> gateway server -> backend api request -> backend api response -> handler callback(caching) -> gateway respojse
func ProxyHTTP(method string, url string, errmsg string, w http.ResponseWriter, r *http.Request, handler respHandlerFunc) {
	req, err := createHTTPRequest(method, url, r)
	if handleHTTPError(errmsg, err, w) {
		return
	}

	resp, err := httpClient.Do(req)
	if handleHTTPError(errmsg, err, w) {
		return
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if handleHTTPError(errmsg, err, w) {
		return
	}

	if resp.StatusCode == http.StatusOK {
		handler(bs)
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(w, bytes.NewReader(bs)); err != nil {
		log.Printf("%s, copy resp failed: %+v", errmsg, err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// handleHTTPError will check http error, log and write response back to client
// if true - error has been found and handled
func handleHTTPError(msg string, err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	if err == missingAuthToken {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("api auth token autopilotapikey must be provided"))
		return true
	}

	log.Printf("%s %+v\n", msg, err)
	w.WriteHeader(http.StatusInternalServerError)

	return true
}

func createHTTPRequest(method string, url string, oreq *http.Request) (*http.Request, error) {
	req, err := http.NewRequest(method, url, oreq.Body)
	if err != nil {
		return nil, err
	}

	token := oreq.Header.Get(apiAuthHeader)
	if token == "" {
		return nil, missingAuthToken
	}

	req.Header.Set(apiAuthHeader, token)
	req.Header.Set("Content-Type", "application/json")
	// copyHeader(req.Header, oreq.Header)
	return req, nil
}
