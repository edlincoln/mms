package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/edlincoln/mms/internal/configs"
	"github.com/edlincoln/mms/pkg/logger"
)

func main() {
	ctx := context.Background()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for _, c := range configs.GetConfigs() {
		err := c.Init(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("Error on init configuration %s", reflect.TypeOf(c)), err)
		}
	}

	// waiting to a sig channel
	sig := <-sigchan
	logger.Info(fmt.Sprintf("Caught signal %v: terminating", sig))

	for _, c := range configs.GetConfigs() {
		err := c.Close(ctx)
		if err != nil {
			logger.Error(fmt.Sprintf("Error on closing configuration %s : ", reflect.TypeOf(c)), err)
		}
	}
}
