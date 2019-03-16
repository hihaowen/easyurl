package main

import (
	"easyurl/config"
	"easyurl/infra/cache/redis"
	"easyurl/infra/db/mysql"
	"html/template"
	"log"
	"net/http"
	sq "github.com/Masterminds/squirrel"
	"time"
)

const rootPath = "/Users/wenzg/workspace/go/src/easyurl"

func init() {
	// 连接数据库
	mysql.Connect(config.MySQLDSN)

	// 连接cache
	redis.Connect(redis.RedisConf{Addr: config.RedisAddr, Pwd: config.RedisPwd})
}

func main() {
	http.HandleFunc("/admin/dev/key", adminDevKeyAddHandler)
	http.HandleFunc("/admin/dev/key/save", adminDevKeySaveHandler)

	log.Fatal(http.ListenAndServe(":8888", nil))
}

// 缓存
var templates = template.Must(template.ParseFiles(rootPath + "/tpl/admin/dev_key_add.html"))

// api_dev_key的创建
func adminDevKeyAddHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "dev_key_add", nil)
}

func adminDevKeySaveHandler(w http.ResponseWriter, r *http.Request) {
	apiDevKey := r.FormValue("api_dev_key")
	userId := r.FormValue("user_id")
	nowTs := time.Now().Unix()

	_, err := sq.
		Insert("api_dev_keys").Columns("api_dev_key", "user_id","create_ts","update_ts").
		Values(apiDevKey, userId, nowTs, nowTs).
		RunWith(mysql.Db).
		Exec()

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin/dev/key", http.StatusFound)
}

// 模版渲染
func renderTemplate(w http.ResponseWriter, r *http.Request, tpl string, page interface{}) {
	err := templates.ExecuteTemplate(w, tpl+".html", page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
