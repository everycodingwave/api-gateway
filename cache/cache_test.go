package cache

import (
	"os"
	"testing"
	"time"

	"github.com/everycodingwave/api-gateway/env"
)

var cacheAddr string

func init() {
	cacheAddr = os.Getenv(env.CacheAddr)
}

func TestBasicOperation(t *testing.T) {
	if cacheAddr == "" {
		t.Logf("skip cache basic test cause %s env is not setting", env.CacheAddr)
		return
	}

	cac := New(cacheAddr)
	err := cac.Set("k1", "v1", 0)
	if err != nil {
		t.Errorf("set k1 v1 error %+v", err)
	}

	val, err := cac.Get("k1")
	if err != nil || val != "v1" {
		t.Errorf("get k1 expect return v1, but return %s %v", val, err)
	}

	val, err = cac.Get("k2")
	if err != KeyNotExisted {
		t.Errorf("get non-existed key k2 should return not existed error, but %s %v", val, err)
	}

	err = cac.Del("k1")
	if err != nil {
		t.Errorf("delete k1 error %v", err)
	}

	val, err = cac.Get("k1")
	if err != KeyNotExisted {
		t.Errorf("get k1 after deleting it should return key not existed error but %s %v", val, err)
	}

}

func TestDelEmptyKey(t *testing.T) {
	if cacheAddr == "" {
		t.Logf("skip cache basic test cause %s env is not setting", env.CacheAddr)
		return
	}

	cac := New(cacheAddr)
	err := cac.Del("random key 123")
	if err != nil {
		t.Errorf("expect delete empty key return no error, but %v", err)
	}
}

func TestGetExpiredKey(t *testing.T) {
	if cacheAddr == "" {
		t.Logf("skip cache basic test cause %s env is not setting", env.CacheAddr)
		return
	}

	cac := New(cacheAddr)
	err := cac.Set("e1", "v1", 0)
	if err != nil {
		t.Errorf("set k1 v1 error %+v", err)
	}

	err = cac.Set("e2", "v2", time.Second*3)
	if err != nil {
		t.Errorf("set e2 v2 error %+v", err)
	}

	val, err := cac.Get("e2")
	if err != nil || val != "v2" {
		t.Errorf("get key e2 should return v2, but %s %v", val, err)
	}

	time.Sleep(time.Second * 4)

	val, err = cac.Get("e2")
	if err != KeyNotExisted {
		t.Errorf("get expired key e2 should return not existed error, but %s %v", val, err)
	}

	val, err = cac.Get("e1")
	if err != nil || val != "v1" {
		t.Errorf("get e1 expect return v1, but return %s %v", val, err)
	}
}
