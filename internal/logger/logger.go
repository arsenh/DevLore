package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(log.InfoLevel)
}

func L() *log.Logger {
	return log.StandardLogger()
}

func LogAndErr(msg string, args ...interface{}) error {
	L().Errorf(msg, args)
	return fmt.Errorf(msg, args)
}
