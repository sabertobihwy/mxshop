package initialize

import "go.uber.org/zap"

func InitilizeLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
