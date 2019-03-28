package api

import (
	"crypto/md5"
	"easyurl/infra/db/mysql"
	"easyurl/model"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	sq "github.com/Masterminds/squirrel"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	apiDevKey := r.FormValue("api_dev_key")
	originalUrl := r.FormValue("original_url")
	// expireDate := time.Now().Unix()
	customAlias := r.FormValue("custom_alias")

	// 校验api_dev_key
	apiDevKeyItem, err := model.GetOneByApiKey(apiDevKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "api_dev_key错误", http.StatusInternalServerError)
	}

	// 自定义 or 自动生成alias
	var urlAlias string
	if customAlias == "" {
		urlAlias = randomAliasGenerator(apiDevKeyItem.ApiDevKey)
	} else {
		urlAlias = customAlias
	}

	// 检查hash是否存在，存在且过期则更新，不存在直接创建
	urlItem, err := model.GetOneByHash(urlAlias)
	if err != nil {
		log.Println(err)
		http.Error(w, "生成失败,请稍后重试", http.StatusInternalServerError)
		return
	}

	nowTs := time.Now().Unix()

	if urlItem.Hash == "" {
		_, err := sq.
			Insert("urls").Columns("hash", "original_url", "user_id", "expire_ts").
			Values(urlAlias, originalUrl, apiDevKeyItem.UserId, nowTs+10).
			RunWith(mysql.Db).
			Exec()

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// 存在但未过期，不允许更改
		if urlItem.ExpireTs > uint64(nowTs) {
			log.Println("创建失败，请重试")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := sq.Update("urls").
			Set("user_id", apiDevKeyItem.UserId).
			Set("original_url", originalUrl).
			Set("expire_ts", nowTs+10).
			RunWith(mysql.Db)

		_, err := b.Exec()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	println(urlAlias)
}

// url密钥算法: substr(base64encode(md5(api_dev_key+sec)), 0, 6)
func randomAliasGenerator(apiDevKey string) string {
	// 原始字符串
	ori := apiDevKey + strconv.FormatInt(time.Now().Unix(), 10)

	// md5
	h := md5.New()
	_, _ = io.WriteString(h, ori)

	md5Str := fmt.Sprintf("%x", h.Sum(nil))

	// base64
	base64Str := base64.StdEncoding.EncodeToString([]byte(md5Str))

	return string([]byte(base64Str)[:6])
}
