package controllers

import (
	"errors"
	"net/http"

	"github.com/edlincoln/mms/internal/http/handlers"
	"github.com/edlincoln/mms/internal/http/request"
	"github.com/edlincoln/mms/internal/http/response"
	"github.com/edlincoln/mms/internal/service"
	"github.com/edlincoln/mms/internal/utils"
	"github.com/edlincoln/mms/internal/utils/constants"
	"github.com/go-chi/chi"
)

type MmsPair struct {
	requestHandler  handlers.RequestHandler
	responseHandler handlers.ResponseHandler
	service         service.MmsPairService
}

func NewMmsPair(requestHandler handlers.RequestHandler, responseHandler handlers.ResponseHandler, service service.MmsPairService) MmsPair {
	return MmsPair{
		requestHandler:  requestHandler,
		responseHandler: responseHandler,
		service:         service,
	}
}

func (c MmsPair) Route(r chi.Router) {
	r.Route("/{pair}", func(idRoute chi.Router) {
		idRoute.Get("/mms", c.findByTimestamp)
	})
}

func (c MmsPair) GetPath() string {
	return constants.MmsPairPath
}

// @Title findMmsPair
// @Tags MmsPair
// @Summary find mms by pair and time stamp
// @Description This resource performs find mms by pair and time stamp
// @Description **Errors http codes response:**
// @Description HTTP | Description | Code | Note
// @Description -----|-----|-----|-----
// @Description 500 | Internal server error | 500 | N/A
// @Param 		pair 	path		string 		true 	"pair identifier"
// @Param   	from 	query    	int     	false 	"timestamp from"
// @Param   	to 		query    	int     	false 	"timestamp to"
// @Param   	range 	query    	int     	false 	"range"
// @Success 	200 	{{object}}  response.MmsPair
// @Failure 	400 	"Bad request"
// @Failure 	500 	"Internal Server Error"
// @Accept json
// @Router /v1/{pair}/mms [get]
func (c MmsPair) findByTimestamp(w http.ResponseWriter, r *http.Request) {

	pair := c.requestHandler.GetURLParam(r, constants.Pair)

	var input request.MmsPairDto
	err := c.requestHandler.BindQuery(r, &input)
	if err != nil {
		c.responseHandler.Exception(w, err)
		return
	}

	// todo: descomentar !!!
	// err = input.Valid()
	// if err != nil {
	// 	c.responseHandler.Exception(w, err)
	// 	return
	// }

	ret, errSave := c.service.FindByPairAndTimestampRange(r.Context(), pair, input)
	if errSave != nil {
		if errors.Is(utils.InternalServerError, errSave) {
			c.responseHandler.Exception(w, handlers.NewInternalServerError(errSave.Error(), errSave))
			return
		}
		c.responseHandler.Exception(w, handlers.NewUnprocessableEntityError(errSave.Error()))
		return
	}

	c.responseHandler.Ok(w, response.ToMmsPairListResponse(input.Range, ret))
}
