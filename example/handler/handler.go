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
