package api

import (
	"easyurl/controller"
	"easyurl/infra/engine"
	"os"
	"time"
)

func GetPid(c *engine.Context) {

	controller.Response(c, 0, "", map[string]interface{}{
		"now_time": time.Now(),
		"now_pid": os.Getpid(),
	})

	return
}
