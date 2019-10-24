package main

import (
	"github.com/autopilot/apigateway/cache"
	"github.com/autopilot/apigateway/env"
	"github.com/autopilot/apigateway/server"
	"log"
	"os"
)

func main() {
	//envs supposed to be stored in some dedicated secret service
	cacheAddr := os.Getenv(env.CacheAddr)
	if cacheAddr == "" {
		log.Fatalf("CACHE_SERVER_ADDR env is missing")
	}

	srv := server.New(cache.New(cacheAddr), server.ProxyHTTP)
	err := srv.Start()
	if err != nil {
		log.Fatalf("server stopped unexpectedly, err %+v\n", err)
	}
}
