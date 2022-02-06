package logutil

import "log"

type MyLogger struct {
	Lgr       log.Logger
	Log_level int
}

func (l *MyLogger) WriteLog(message string, log_level int) {
	if log_level <= l.Log_level {
		l.Lgr.Printf("[%v] %v\n", log_level, message)
	}
}
