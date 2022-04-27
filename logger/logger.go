package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	logTimeFormat = "2006-01-02 15:04:05"
)

type Config struct {
	Level   logrus.Level
	OutJson bool
	NoColor bool
}

func Init(name string, c Config) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     !c.NoColor,
		TimestampFormat: logTimeFormat,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: fmt.Sprintf("[%s] time", name),
		},
	},
	)

	logrus.SetLevel(c.Level)
	logrus.SetOutput(os.Stdout)

	if c.Level > logrus.DebugLevel {
		logrus.SetReportCaller(true)
	}

	if c.OutJson {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: logTimeFormat,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: fmt.Sprintf("[%s] time", name),
			},
		},
		)
	}
}
