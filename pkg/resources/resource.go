package resources

import (
	"embed"
)

type ResourceImpl struct{}

type Resource interface {
	ReadFile(resource embed.FS, resourcePath string) ([]byte, error)
}

func NewResource() Resource {
	return ResourceImpl{}
}

func (r ResourceImpl) ReadFile(resource embed.FS, resourcePath string) ([]byte, error) {
	appContent, errResources := resource.ReadFile(resourcePath)
	if errResources != nil {
		return nil, errResources
	}
	return appContent, nil
}
