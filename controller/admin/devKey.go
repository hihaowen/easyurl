package admin

import (
	"easyurl/infra/db/mysql"
	sq "github.com/Masterminds/squirrel"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// 获取当前目录地址
var rootPath = os.Getenv("GOPATH") + "/src/easyurl"

// 缓存
var templates = template.Must(template.ParseFiles(rootPath + "/tpl/admin/dev_key_add.html"))

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	apiDevKey := r.FormValue("api_dev_key")
	userId := r.FormValue("user_id")
	nowTs := time.Now().Unix()

	res, err := sq.
		Insert("api_dev_keys").Columns("api_dev_key", "user_id", "create_ts", "update_ts").
		Values(apiDevKey, userId, nowTs, nowTs).
		RunWith(mysql.Db).
		Exec()

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println(res)

	http.Redirect(w, r, "/admin/dev/key?", http.StatusFound)
}

// api_dev_key的创建
func AddHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "dev_key_add", nil)
}

// 模版渲染
func renderTemplate(w http.ResponseWriter, r *http.Request, tpl string, page interface{}) {
	err := templates.ExecuteTemplate(w, tpl+".html", page)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
