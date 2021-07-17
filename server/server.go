package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"omarkhd/memkv/store"
)

type server struct {
	storage store.Store
}

func New(storage store.Store) (*server, error) {
	srv := &server{storage: storage}
	http.HandleFunc("/", srv.handle)
	return srv, nil
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	// All requests should be prefixed
	if !strings.HasPrefix(r.URL.Path, "/keys") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// If no storage no data
	if s.storage == nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	// List all keys
	if r.URL.Path == "/keys" {
		s.ls(w)
		return
	}
	// Getting key from path
	parts := strings.SplitN(r.URL.Path, "/", 3)
	if len(parts) != 3 || parts[2] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// What to do with the key?
	key := parts[2]
	switch r.Method {
	case http.MethodGet:
		s.get(w, key)
	case http.MethodPut, http.MethodPost:
		s.put(w, r, key)
	case http.MethodDelete:
		s.delete(w, key)
	}
}

func (s *server) ls(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	if err := encoder.Encode(s.storage.Keys()); err != nil {
		log.Printf(err.Error())
	}
}

func (s *server) get(w http.ResponseWriter, key string) {
	value := s.storage.Get(key)
	if _, err := io.WriteString(w, value); err != nil {
		log.Printf(err.Error())
	}
}

func (s *server) put(w http.ResponseWriter, r *http.Request, key string) {
	value, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	s.storage.Put(key, string(value))
}

func (s *server) delete(w http.ResponseWriter, key string) {
	s.storage.Delete(key)
}

func (s *server) Start() {
	log.Fatal(http.ListenAndServe(":4444", nil))
}
