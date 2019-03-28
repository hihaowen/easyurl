package api

import (
	"easyurl/controller"
	"net/http"
	"time"
)

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	controller.Response(w, 0, "", map[string]int64{"now_time": time.Now().Unix()})
}
