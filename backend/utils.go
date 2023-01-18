package main

import "prc_hub_back/domain/model/logger"

func convertLogLevel(str string) (ok bool, lv logger.Level) {
	switch str {
	case "panic":
		return true, logger.PanicLevel
	case "fatal":
		return true, logger.FatalLevel
	case "error":
		return true, logger.ErrorLevel
	case "warn":
		return true, logger.WarnLevel
	case "info":
		return true, logger.InfoLevel
	case "debug":
		return true, logger.DebugLevel
	case "trace":
		return true, logger.TraceLevel
	default:
		return false, logger.InfoLevel
	}
}
