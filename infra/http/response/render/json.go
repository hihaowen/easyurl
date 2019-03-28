package render

import (
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
	WriteContentType(w, contentType)
}

func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	WriteContentType(w, contentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
