package server

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gorilla/mux"

	"github.com/everycodingwave/api-gateway/env"
)

func (s *apiServer) getRouter() *mux.Router {
	// for this task this way might be messy, but for real application writing this way keeps things cleaner
	routes := map[string]map[string]http.HandlerFunc{
		"GET": {
			"/v1/contact/{contact_id}": s.getContact,
		},

		"POST": {
			"/v1/contact": s.createContact,
		},

		"PUT": {
			"/v1/contact": s.createContact,
		},
	}

	router := mux.NewRouter()
	for method, pathMap := range routes {
		for path, handler := range pathMap {
			router.HandleFunc(path, handler).Methods(method)
		}
	}

	return router
}

func (s *apiServer) getContact(w http.ResponseWriter, r *http.Request) {
	contactID, ok := mux.Vars(r)["contact_id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("contact_id is missing\r\n"))
		return
	}

	obj, err := s.cac.Get(contactID)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		if _, err = io.Copy(w, bytes.NewReader([]byte(obj))); err != nil {
			log.Printf("[apiserver error]get contact: copy cache content to resp failed, contact id %s, err %+v", contactID, err)
		}

		return
	}

	s.proxyFunc("GET", contactAPIURL+"/"+contactID, "[apiserver error]get contact error:", w, r, func(bs []byte) {
		// this callback will be called after getting response of the backend api server succussfully
		expireSec := 0

		v, err := strconv.Atoi(os.Getenv(env.CacheExpiredSec))
		if err == nil {
			expireSec = v
		}

		err = s.cac.Set(contactID, string(bs), time.Duration(expireSec))
		if err != nil {
			log.Printf("[apiserver error]set cache failed, contactID %s err %v", contactID, err)
		}
	})
}

// createContact may create or update contact based on whether user's email is already there.
// since the backend handle all the logic, this gateway server will try not get into any business logic
// so as long as calling backend api return 200 OK it will parse out the contract_id and then invalidate cache
// this behavior might waste a bit redis network io if it's just a creating request.
func (s *apiServer) createContact(w http.ResponseWriter, r *http.Request) {
	s.proxyFunc("POST", contactAPIURL, "[apiserver error]create contact error:", w, r, func(bs []byte) {
		contactID, err := jsonparser.GetString(bs, "contact_id")
		if err != nil {
			log.Printf("[apiserver error]createContact api, parse resp failed, contactID %s err %v", contactID, err)
			return
		}

		err = s.cac.Del(contactID)
		if err != nil {
			// purge cache failed can be quite tricky, cause user might have inconsistent data
			// this reply on wether we set the cache expiring time, if that time window is samll
			// the inconsistent data problem can be reduce down a bit, but still there.
			// should keep an eye on this log whenever it occurs, and consider wrting a tool to deal with this
			log.Printf("[apiserver error]purge cache failed, contactID %s %v", contactID, err)
		}

	})
}
