package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

type Logger struct {
	*logrus.Entry
}

var instance Logger

func GetLogger(level string) *Logger {

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	}

	l.SetOutput(os.Stdout)
	l.SetLevel(logLevel)

	instance = Logger{logrus.NewEntry(l)}

	return &instance

}
