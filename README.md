# oapi-extensions-middleware


OpenAPI specification supports extensions in paths,server,schema.That allows the users to extend the oapi file for custom use cases. The [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen/) doesn't support extensions at path level. This middleware adds the extensions support at path level. The keys which start with `x-` will be added to the request context. This can come in handy for many usecases like selective upgrading requests to sockets based on spec or selective skip authentication by adding `x-skip-auth:true` and then using the value of the key in the middleware to skip authentication etc.

Note: Because of Strongly typed nature of Golang the values are available as []byte in the echo request context, One can use json.Unmarshal to unmarshal the []byte to corresponding types.


Example:

> contract.yaml
```yaml 
openapi: "3.0.0"
info:
  version: 1.0.0
  title: oapi-extensions middleware
paths:
  /test:
    get:
      x-test: this is test
      description: test endpoint
      operationId: test
      responses:
        '200':
          description: success response
          content:
            application/json:
              schema:
                type: object
```

Command to generate 
> oapi-codegen -o ./gen/example.gen.go --package gen ./contract.yaml 


> main.go

```go
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

	handler := handler.Handler{}

	gen.RegisterHandlers(server, &handler)

	server.Start(":1232")
}
```

handler/handler.go
```go
package handler

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Handler struct {
}

func (h *Handler) Test(ctx echo.Context) error {
	value := ctx.Get("x-test")

	valByte := value.([]byte)

	var val string

	json.Unmarshal(valByte, &val)

	log.Info("Value of x-test:", val)

	return ctx.String(200, val)
}
```

To Run the example:

> cd example

If you changed the spec run
> oapi-codegen -o ./gen/example.gen.go --package gen ./contract.yaml 

> go run main.go

In a new terminal run 

> curl localhost:1232/test
