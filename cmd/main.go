package main

import (
	"github.com/everycodingwave/api-gateway/cache"
	"github.com/everycodingwave/api-gateway/env"
	"github.com/everycodingwave/api-gateway/server"
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
