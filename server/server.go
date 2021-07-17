package server

import (
	"log"
	"net/http"
)

type server struct {
}

func New() (*server, error) {
	s := &server{}
	http.HandleFunc("/", s.handle)
	return s, nil
}

func (s *server) handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *server) Start() {
	log.Fatal(http.ListenAndServe(":4444", nil))
}
