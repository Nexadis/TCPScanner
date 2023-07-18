package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func init() {
	log, _ := zap.NewDevelopment()
	Log = log.Sugar()
}
