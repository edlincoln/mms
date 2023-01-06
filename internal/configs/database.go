package configs

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/edlincoln/mms/internal/container"
	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/pkg/logger"
	_ "github.com/lib/pq"
)

type DatabaseConfiguration struct {
	client *sql.DB
}

func (c *DatabaseConfiguration) Init(ctx context.Context) error {
	logger.Info("Starting database")

	var err error
	client, err := sql.Open("postgres", c.getConnectionConfig())
	if err != nil {
		panic(err)
	}

	if err = client.Ping(); err != nil {
		panic(err)
	}
	c.client = client
	container.Container().DatabaseManager.SetClient(client)

	log.Println("database successfully configured")

	return nil
}

func (c *DatabaseConfiguration) Close(ctx context.Context) error {
	return c.client.Close()
}

func (c *DatabaseConfiguration) getConnectionConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		global.AppConfig.Database.Connection.Host,
		global.AppConfig.Database.Connection.Port,
		global.AppConfig.Database.Connection.User,
		global.AppConfig.Database.Connection.Password,
		global.AppConfig.Database.Connection.DatabaseName,
		global.AppConfig.Database.Connection.SSLMode)
}
