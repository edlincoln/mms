package controllers

import (
	"net/http"

	"github.com/edlincoln/mms/internal/http/handlers"
	"github.com/edlincoln/mms/internal/service"
	"github.com/edlincoln/mms/internal/utils/constants"
	"github.com/go-chi/chi"
)

type Extractor struct {
	requestHandler  handlers.RequestHandler
	responseHandler handlers.ResponseHandler
	service         service.ExtractorManagerService
}

func NewExtractor(requestHandler handlers.RequestHandler, responseHandler handlers.ResponseHandler, service service.ExtractorManagerService) Extractor {
	return Extractor{
		requestHandler:  requestHandler,
		responseHandler: responseHandler,
		service:         service,
	}
}

func (c Extractor) Route(r chi.Router) {
	r.Route("/daily", func(idRoute chi.Router) {
		idRoute.Post("/", c.daily)
	})
}

func (c Extractor) GetPath() string {
	return constants.Extractor
}

// @Title dailyExtractor
// @Tags Extractor
// @Summary find mms by pair and time stamp
// @Description This resource performs find mms by pair and time stamp
// @Description **Errors http codes response:**
// @Description HTTP | Description | Code | Note
// @Description -----|-----|-----|-----
// @Description 500 | Internal server error | 500 | N/A
// @Success 	201
// @Failure 	500 	"Internal Server Error"
// @Accept json
// @Router /v1/extractor/daily [post]
func (c Extractor) daily(w http.ResponseWriter, r *http.Request) {
	err := c.service.LoadDailyData(r.Context(), 0)
	if err != nil {
		c.responseHandler.Exception(w, handlers.NewInternalServerError(err.Error(), err))
		return
	}

	c.responseHandler.Created(w, nil)
}
