// 获取相关配置
package config

import (
	"os"
	"strconv"
	"strings"
)

// getString is a string environment variable parsing and
// default-setting function
func getString(name string, def interface{}) string {
	val := os.Getenv(name)
	// Make sure certain values have must sets
	if val == "" && def == nil {
		panic(name)
	}

	if val != "" {
		return val
	}
	return def.(string)
}

// getBool is a bool environment variable parsing and
// default-setting function
func getBool(name string, def bool) bool {
	val := os.Getenv(name)
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return ret
}

// getInt is an int environment variable parsing and
// default-setting function
func getInt(name string, def int64) int64 {
	val := os.Getenv(name)
	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return def
	}
	return ret
}

// getStringArray is a string array (separated by ",") environment variable
// parsing and default-setting function
func getStringArray(name string, def []string) []string {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	// Split into a slice
	return strings.Split(val, ",")
}

// Db
var (
	MySQLDSN  = getString("MySQL_DSN", "root:root@tcp(127.0.0.1:3306)/easyurl")
)

// Cache
var (
	RedisAddr = getString("Redis_Addr", "127.0.0.1:6379")
	RedisPwd = getString("Redis_Pwd", "")
)

const SAAS = "saas"
