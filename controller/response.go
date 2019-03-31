package controller

import (
	"easyurl/infra/engine"
	"easyurl/infra/http/response/render"
	"net/http"
)

type r map[string]interface{}

func Response(c *engine.Context, errno int32, error string, data interface{}) {
	msg := &r{
		"errno": errno,
		"error": error,
		"data":  data,
	}

	c.Render(http.StatusOK, render.JSON{Data: msg})

	return
}
