package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pingpulse/internal/monitor"
)

type Server struct {
	sink   interface {
		All() map[string][]monitor.Result
		Latest(url string, n int) []monitor.Result
	}
	addr string
	srv  *http.Server
}

func NewServer(sink interface {
	All() map[string][]monitor.Result
	Latest(url string, n int) []monitor.Result
}, addr string) *Server {
	mux := http.NewServeMux()
	s := &Server{sink: sink, addr: addr}
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/results", s.handleResults)
	s.srv = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return s
}

func (s *Server) Start() error {
	log.Printf("http listening on %s", s.addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if url == "" {
		_ = enc.Encode(s.sink.All())
		return
	}
	_ = enc.Encode(s.sink.Latest(url, 20))
}