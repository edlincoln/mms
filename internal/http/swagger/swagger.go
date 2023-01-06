package swagger

import (
	"fmt"
	"net/http"

	"github.com/edlincoln/mms/internal/global"
	"github.com/edlincoln/mms/internal/http/swagger/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

func ConfigurationSwagger() http.HandlerFunc {

	docs.SwaggerInfo.Title = global.AppConfig.App.Name
	docs.SwaggerInfo.Version = global.AppConfig.App.Version
	docs.SwaggerInfo.Description = global.AppConfig.App.Description
	docs.SwaggerInfo.BasePath = global.AppConfig.Swagger.BasePath
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", global.AppConfig.Swagger.Ip, global.AppConfig.Swagger.Port)

	return httpSwagger.Handler(httpSwagger.URL(global.AppConfig.Swagger.Custom.BasePath))
}

// ReadSwaggerDoc :
func ReadSwaggerDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		s, _ := swag.ReadDoc()
		_, _ = w.Write([]byte(s))
	}
}
