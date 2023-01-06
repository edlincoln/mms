package configs

import (
	"context"

	embed "github.com/edlincoln/mms"
	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/pkg/logger"
	"github.com/edlincoln/mms/pkg/resources"
	"github.com/edlincoln/mms/pkg/viper"
)

const (
	_applicationFileName = "application"
	_extension           = "yml"
	_resourcePathViper   = "resources/application.yml"
)

type ViperConfiguration struct {
}

func (c *ViperConfiguration) Init(context.Context) error {

	application := global.Application{}

	app := viper.NewApplication(_resourcePathViper, _applicationFileName, _extension, viper.NewIViper(), resources.NewResource())

	_, err := app.ConfigurationViper(&application, embed.Resources)
	if err != nil {
		logger.Error("error initializing properties file", err)
	}

	global.AppConfig = application

	return nil
}

func (c *ViperConfiguration) Close(_ context.Context) error {
	return nil
}
