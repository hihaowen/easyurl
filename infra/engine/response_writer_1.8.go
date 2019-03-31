package engine

import "net/http"

// ResponseWriter ...
type ResponseWriter interface {
	responseWriterBase
	// get the http.Pusher for server push
	Pusher() http.Pusher
}

func (w *responseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
