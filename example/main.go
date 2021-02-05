package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/yolossn/oapi-extensions-middleware/example/gen"
	"github.com/yolossn/oapi-extensions-middleware/example/handler"
	"github.com/yolossn/oapi-extensions-middleware/middleware"
)

func main() {

	server := echo.New()

	swagger, err := gen.GetSwagger()
	if err != nil {
		log.Fatal("Couldn't read swagger spec", err)
	}

	server.Use(middleware.OapiAddPathExtensions(swagger))

	// Using yaml file vvv
	// ---------

	// pathMiddleware, err := middleware.OapiAddPathExtensionsFromYamlFile("./contract.yaml")
	// if err != nil {
	// 	log.Fatal("Couldn't read swagger spec", err)
	// }

	// server.Use(pathMiddleware)
	// ----------

	handler := handler.Handler{}

	gen.RegisterHandlers(server, &handler)

	server.Start(":1232")
}
