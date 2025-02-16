package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer(host string, port string, timeout, idletimeout time.Duration) *Server {
	srv := &http.Server{
		Addr:           host + ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		IdleTimeout:    idletimeout,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer.Handler = handler
	return s.httpServer.ListenAndServe()
}

func (s *Server) Gracefull(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
