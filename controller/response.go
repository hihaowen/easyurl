package controller

import (
	"easyurl/infra/http/response/render"
	"net/http"
)

type r map[string]interface{}

func Response(w http.ResponseWriter, errno int32, error string, data interface{}) {
	res := &r{
		"errno": errno,
		"error": error,
		"data":  data,
	}

	if err := (render.JSON{Data: res}.Render(w)); err != nil {
		panic(err)
	}
}
