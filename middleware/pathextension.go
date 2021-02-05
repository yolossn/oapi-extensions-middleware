package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// OapiAddPathExtensionsFromYamlFile returns OapiPathExtensions middleware from yaml file
func OapiAddPathExtensionsFromYamlFile(path string) (echo.MiddlewareFunc, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", path, err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s as Swagger YAML: %s",
			path, err)
	}
	return OapiAddPathExtensions(swagger), nil
}

// OapiAddPathExtensions returns OapiExtensionsMiddleware from swagger spec
func OapiAddPathExtensions(swagger *openapi3.Swagger) echo.MiddlewareFunc {
	return OapiAddPathExtensionsWithOptions(swagger, nil)
}

// OapiAddPathExtensionsWithOptions returns OapiExtensionsMiddleware from swagger spec and custom options
func OapiAddPathExtensionsWithOptions(swagger *openapi3.Swagger, options *middleware.Options) echo.MiddlewareFunc {
	router := openapi3filter.NewRouter().WithSwagger(swagger)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			// Skipper logic
			skipper := getSkipperFromOptions(options)
			if skipper(ctx) {
				return next(ctx)
			}
			req := ctx.Request()

			// Find Route
			route, _, err := router.FindRoute(req.Method, req.URL)
			if err != nil {
				return next(ctx)
			}

			// Fetch extensions of operation and set it as string in context
			operation := route.PathItem.GetOperation(req.Method)

			for key, value := range operation.Extensions {
				if value != nil {
					if val, ok := value.(json.RawMessage); ok {
						out, err := val.MarshalJSON()
						if err != nil {
							continue
						}
						ctx.Set(key, out)
					}
				}
			}

			return next(ctx)
		}
	}
}

func getSkipperFromOptions(options *middleware.Options) echomiddleware.Skipper {
	if options == nil {
		return echomiddleware.DefaultSkipper
	}

	if options.Skipper == nil {
		return echomiddleware.DefaultSkipper
	}

	return options.Skipper
}
