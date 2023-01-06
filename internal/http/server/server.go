package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/internal/http/server/routers"
	"github.com/edlincoln/mms/pkg/logger"
)

type HttpServer struct {
	start  time.Time
	server *http.Server
}

func NewHttpServer(start time.Time) *HttpServer {
	return &HttpServer{start: start}
}

func (s *HttpServer) StartHttpServer(ctx context.Context) {

	logger.Info("Starting server HTTP")

	r := routers.NewRouter(ctx)
	logger.Info("Created routers")

	port := global.AppConfig.Server.Port
	logger.Info(fmt.Sprintf("Configured server port: %d", port))
	s.server = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r.Router}
	go func() {
		e := s.server.ListenAndServe()
		if e != nil {
			if strings.Contains(e.Error(), "Server closed") {
				logger.Info("Stopped HTTP Server")
				return
			}
			logger.Error("Internal error", e)
		}
	}()
}

func (s HttpServer) Stop(ctx context.Context) error {
	if s.server != nil {
		logger.Info("Stopping HTTP Server")
		return s.server.Shutdown(ctx)
	}
	return nil
}
