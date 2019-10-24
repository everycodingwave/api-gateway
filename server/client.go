package server

import (
	"net"
	"net/http"
	"time"
)

const (
	httpConnectTimeout time.Duration = 5 * time.Second
	httpIOTimeout      time.Duration = 10 * time.Second
)

var transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: httpConnectTimeout,
	}).Dial,
	TLSHandshakeTimeout: httpConnectTimeout,
}
var httpClient = &http.Client{
	Timeout:   httpIOTimeout,
	Transport: transport,
}
