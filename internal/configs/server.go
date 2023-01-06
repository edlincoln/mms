package configs

import (
	"context"
	"time"

	"github.com/edlincoln/mms/internal/http/server"
	"github.com/edlincoln/mms/pkg/logger"
)

type ServerConfiguration struct {
	httpServer *server.HttpServer
}

func (c *ServerConfiguration) Init(ctx context.Context) error {
	c.httpServer = server.NewHttpServer(time.Now())
	c.httpServer.StartHttpServer(ctx)
	return nil
}

func (c *ServerConfiguration) Close(ctx context.Context) error {
	logger.Info("Closing Server...")
	return c.httpServer.Stop(ctx)
}
