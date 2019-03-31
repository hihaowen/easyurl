package main

import (
	"easyurl/config"
	api2 "easyurl/controller/api"
	"easyurl/infra/cache/redis"
	"easyurl/infra/db/mysql"
	"easyurl/infra/engine"
	"log"
)

func init() {
	// 连接数据库
	mysql.Connect(config.MySQLDSN)
	// 连接cache
	redis.Connect(redis.RedisConf{Addr: config.RedisAddr, Pwd: config.RedisPwd})
}

func Limit() engine.HandlerFunc {
	return func(c *engine.Context) {
		log.Printf("%s %s %+v\n", c.ClientIP(), c.Request.URL, c.Request.Header)
		c.Next()
	}
}

func main() {
	/*
	http.HandleFunc("/admin/dev/key", admin.AddHandler)
	http.HandleFunc("/admin/dev/key/save", admin.SaveHandler)
	http.HandleFunc("/api/url/create", api.CreateHandler)

	log.Fatal(http.ListenAndServe(":8888", nil))
	*/
	r := engine.New()
	r.Use(Limit())

	v1 := r.Group("/v1")
	{
		api := v1.Group("/api")
		{
			api.GET("/pid", api2.GetPid)
		}
	}

	_ = r.Run(":8888")
}
