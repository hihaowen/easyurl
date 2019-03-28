package render

import (
	"easyurl/infra/http/response"
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data interface{}
}

var contentType = []string{"application/json; charset=utf-8"}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	response.WriteContentType(w, contentType)
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	response.WriteContentType(w, contentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
