package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

var responseHandlers ResponseHandler
var onceResponse sync.Once

type ResponseHandler interface {
	Exception(w http.ResponseWriter, err RestErr)
	Created(w http.ResponseWriter, data interface{})
	Ok(w http.ResponseWriter, data interface{})
}

type ResponseHandlerImpl struct{}

func NewResponseHandlers() ResponseHandler {
	return ResponseHandlerImpl{}
}

func GetInstanceResponseHandler() ResponseHandler {
	onceResponse.Do(func() {
		responseHandlers = NewResponseHandlers()
	})
	return responseHandlers
}

func (h ResponseHandlerImpl) Ok(w http.ResponseWriter, data interface{}) {
	response(w, data, http.StatusOK)
}

func (h ResponseHandlerImpl) Created(w http.ResponseWriter, data interface{}) {
	response(w, data, http.StatusCreated)
}

func (h ResponseHandlerImpl) Exception(w http.ResponseWriter, err RestErr) {
	response(w, err, err.Status())
}

func headers(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func response(w http.ResponseWriter, data interface{}, httpStatus int) {
	headers(w)
	w.WriteHeader(httpStatus)
	if data != nil {
		if bytes, e := json.Marshal(data); e != nil {
			handlers := NewInternalServerError("error on json marshal", e)
			bytes, _ := json.Marshal(handlers)
			_, _ = w.Write(bytes)
		} else {
			_, _ = w.Write(bytes)
		}
	}
}
