package server

import (
	"context"
	"fmt"
	"net/http"
	"ping-pong/internal/config"
	"time"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer(cfg *config.Config) *HTTPServer {
	return &HTTPServer{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf(":%s", cfg.ServerConfig.Port),
			MaxHeaderBytes: cfg.ServerConfig.MaxHeaderMegabytes << 28,
		},
	}
}

func (s *HTTPServer) Run(port string) error {
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
