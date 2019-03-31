package engine

import "os"

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)
const (
	debugCode = iota
	releaseCode
	testCode
)

// 环境模式
var envMode = debugCode
var envModeName = DebugMode

func SetMode(value string) {
	switch value {
	case DebugMode, "":
		envMode = debugCode
	case ReleaseMode:
		envMode = releaseCode
	case TestMode:
		envMode = testCode
	default:
		panic("mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	envModeName = value
}

func init() {
	mode := os.Getenv("ENV_MODE")
	SetMode(mode)
}
