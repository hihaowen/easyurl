package model

import (
	"easyurl/infra/db/mysql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type UrlItem struct {
	Hash        string `json:"hash"`
	OriginalUrl string `json:"original_url"`
	UserId      uint32 `json:"user_id"`
	ExpireTs    uint64 `json:"expire_ts"`
}

// 根据hash获取
func GetOneByHash(hash string) (UrlItem, error) {
	sql, args, err := sq.Select("*").
		From("urls").
		Where(sq.Eq{"hash": hash}).
		Limit(1).
		ToSql()

	if err != nil {
		log.Println(err)
		return UrlItem{}, err
	}

	var urlItem UrlItem
	if row := mysql.Db.QueryRow(sql, args...); row != nil {
		err := row.Scan(&urlItem.Hash, &urlItem.OriginalUrl, &urlItem.UserId, &urlItem.ExpireTs)
		if err != nil {
			return UrlItem{}, nil
		}
	}

	return urlItem, nil
}
