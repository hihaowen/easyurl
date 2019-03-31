package engine

import "log"

func IsDebugging() bool {
	return envMode == debugCode
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		log.Printf("[debug] "+format, values...)
	}
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}
