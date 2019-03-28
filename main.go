package main

import (
	"easyurl/config"
	"easyurl/controller/admin"
	"easyurl/controller/api"
	"easyurl/infra/cache/redis"
	"easyurl/infra/db/mysql"
	"log"
	"net/http"
)

func init() {
	// 连接数据库
	mysql.Connect(config.MySQLDSN)
}

func init() {
	// 连接cache
	redis.Connect(redis.RedisConf{Addr: config.RedisAddr, Pwd: config.RedisPwd})
}

func main() {
	http.HandleFunc("/admin/dev/key", admin.AddHandler)
	http.HandleFunc("/admin/dev/key/save", admin.SaveHandler)
	http.HandleFunc("/api/url/create", api.CreateHandler)
	http.HandleFunc("/api/now/time", api.TimeHandler)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
