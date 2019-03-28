package model

import (
	"easyurl/infra/db/mysql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

type ApiDevKeyItem struct {
	ApiDevKey string `json:"api_dev_key"`
	UserId    uint32 `json:"user_id"`
	Status    uint8  `json:"status"`
	CreateTs  uint64 `json:"create_ts"`
	UpdateTs  uint64 `json:"update_ts"`
}

func GetOneByApiKey(apiDevKey string) (ApiDevKeyItem, error) {
	sql, args, err := sq.Select("*").
		From("api_dev_keys").
		Where(sq.Eq{"status": 2, "api_dev_key": apiDevKey}).
		Limit(1).
		ToSql()

	if err != nil {
		log.Println(err)
		return ApiDevKeyItem{}, err
	}

	var apiDevKeyItem ApiDevKeyItem
	if row := mysql.Db.QueryRow(sql, args...); row != nil {
		err := row.Scan(&apiDevKeyItem.ApiDevKey, &apiDevKeyItem.UserId, &apiDevKeyItem.Status, &apiDevKeyItem.CreateTs, &apiDevKeyItem.UpdateTs)
		if err != nil {
			return ApiDevKeyItem{}, err
		}
	}

	return apiDevKeyItem, nil
}
