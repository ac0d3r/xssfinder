package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	Init(Config{
		Level:   logrus.DebugLevel,
		NoColor: true,
	})
	logrus.Debugln("tessss")
}
