package render

import (
	"net/http"
)

type Render interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

// 这里面起到了接口规范的作用
var (
	_ Render = JSON{}
)

func WriteContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
