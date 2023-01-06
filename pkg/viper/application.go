package viper

import (
	"bytes"
	"embed"

	"github.com/edlincoln/mms/pkg/resources"
	"github.com/spf13/viper"
)

type ApplicationImpl struct {
	resourcePath    string
	applicationName string
	extension       string
	viper           IViper
	resource        resources.Resource
}

func NewApplication(resourcePath, applicationName, extension string, v IViper, resource resources.Resource) Application {
	return ApplicationImpl{
		resourcePath:    resourcePath,
		applicationName: applicationName,
		extension:       extension,
		viper:           v,
		resource:        resource,
	}
}

// ConfigurationViper godoc
func (a ApplicationImpl) ConfigurationViper(model interface{}, resource embed.FS) (*viper.Viper, error) {

	appContent, errResources := a.resource.ReadFile(resource, a.resourcePath)
	if errResources != nil {
		return nil, errResources
	}

	viperConfig, errConfig := a.viper.Configuration(bytes.NewBuffer(appContent), a.extension, a.applicationName)
	if errConfig != nil {
		return nil, errConfig
	}

	errUnmarshal := a.viper.Unmarshal(viperConfig, &model)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return viperConfig, nil
}
