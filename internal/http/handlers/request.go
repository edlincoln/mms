package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var requestHandler RequestHandler
var onceRequest sync.Once

type RequestHandler interface {
	GetURLParam(r *http.Request, key string) string
	BindJson(r *http.Request, destination interface{}) RestErr
	BindQuery(r *http.Request, destination interface{}) RestErr
}

type RequestHandlerImpl struct {
	validator *validator.Validate
}

func NewRequestHandlers() RequestHandler {
	return RequestHandlerImpl{
		validator: initValidator(),
	}
}

func GetInstanceRequestHandler() RequestHandler {
	onceRequest.Do(func() {
		requestHandler = NewRequestHandlers()
	})
	return requestHandler
}

func (h RequestHandlerImpl) GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (h RequestHandlerImpl) BindJson(r *http.Request, destination interface{}) RestErr {
	return h.bind(r, destination)
}

func (h RequestHandlerImpl) BindQuery(r *http.Request, destination interface{}) RestErr {
	return h.query(r, destination)
}

func (h RequestHandlerImpl) bind(r *http.Request, destination interface{}) RestErr {
	if e := render.DecodeJSON(r.Body, destination); e != nil {
		return NewBadRequestError("Invalid json body")
	}
	return h.validateJsonFields(destination)
}

func (h RequestHandlerImpl) query(r *http.Request, destination interface{}) RestErr {

	bytes, errMarshal := json.Marshal(r.URL.Query())
	if errMarshal != nil {
		return NewBadRequestError("error marshaling json")
	}

	newQuery := h.normalizeQuery(bytes)
	if e := json.Unmarshal([]byte(newQuery), destination); e != nil {
		return NewBadRequestError("error marshaling json")
	}

	return h.validateJsonFields(destination)
}

func (h RequestHandlerImpl) normalizeQuery(query []byte) string {
	text := fmt.Sprint(string(query))
	text = strings.Replace(text, "[", "", -1)
	text = strings.Replace(text, "]", "", -1)
	return text
}

func (h RequestHandlerImpl) validateJsonFields(input interface{}) RestErr {
	if err := h.validator.Struct(input); err != nil {
		return NewBadRequestError("Validator error")
	}
	return nil
}

func initValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		alias, found := field.Tag.Lookup("alias")
		if !found {
			return strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		}
		return alias
	})
	return v
}
