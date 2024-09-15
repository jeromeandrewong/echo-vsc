package logger

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	Info  func(msg interface{}, keyvals ...interface{})
	Warn  func(msg interface{}, keyvals ...interface{})
	Error func(msg interface{}, keyvals ...interface{})
	Fatal func(msg interface{}, keyvals ...interface{})
	Debug func(msg interface{}, keyvals ...interface{})
)

func init() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: false,
	})

	Info = logger.Info
	Warn = logger.Warn
	Error = logger.Error
	Fatal = logger.Fatal
	Debug = logger.Debug
}
