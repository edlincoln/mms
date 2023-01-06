package routers

import (
	"sync"

	"github.com/edlincoln/mms/internal/container"
	"github.com/edlincoln/mms/internal/http/controllers"
	"github.com/edlincoln/mms/internal/http/handlers"
)

var controller *ControllerImpl
var once sync.Once

type Controller interface {
	MmsPairController() controllers.MmsPair
	ExtractorController() controllers.Extractor
}

type ControllerImpl struct {
	container *container.ServiceContainer
}

func GetInstance() Controller {
	once.Do(func() {
		controller = &ControllerImpl{
			container: container.Container(),
		}
	})
	return controller
}

func (c ControllerImpl) MmsPairController() controllers.MmsPair {
	return controllers.NewMmsPair(
		handlers.GetInstanceRequestHandler(),
		handlers.GetInstanceResponseHandler(),
		c.container.MmsPairService)
}

func (c ControllerImpl) ExtractorController() controllers.Extractor {
	return controllers.NewExtractor(
		handlers.GetInstanceRequestHandler(),
		handlers.GetInstanceResponseHandler(),
		c.container.ExtractorManagerService)
}
