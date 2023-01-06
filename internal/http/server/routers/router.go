package routers

import (
	"context"
	"fmt"

	"github.com/edlincoln/mms/internal/http/swagger"
	"github.com/edlincoln/mms/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ChiRouter struct {
	Router  chi.Router
	routers []Router
}

func NewRouter(ctx context.Context) ChiRouter {
	chiRouter := ChiRouter{Router: chi.NewRouter()}

	chiRouter.initializeControllers()
	chiRouter.configurationRouters(ctx)

	return chiRouter
}

func (c *ChiRouter) initializeControllers() {
	routes := NewRouterComponent()
	c.routers = []Router{
		routes.MmsPairController,
		routes.ExtractorController,
	}
}

func (c ChiRouter) configurationRouters(ctx context.Context) {
	logger.Info("Creating routers")
	c.Router.Route("/v1", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		for _, router := range c.routers {
			logger.Info(fmt.Sprintf("Router %s created", router.GetPath()))
			r.Route(router.GetPath(), router.Route)
		}
	})

	c.Router.Get("/swagger/*", swagger.ConfigurationSwagger())
	c.Router.Get("/swagger.json", swagger.ReadSwaggerDoc())
}
