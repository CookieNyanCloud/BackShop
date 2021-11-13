package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField( k string, v interface{}) Logger {
	return Logger{l.WithField(k,v)}
}


type writerHook struct {
	Writer []io.Writer
	LogLevels []logrus.Level
}

func (h *writerHook) Fire (entry *logrus.Entry)error {
	line, err := entry.String()
	if err!= nil {
		return err
	}
	for _, w := range h.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (h *writerHook) Levels () []logrus.Level{
	return h.LogLevels
}



func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%S:%d", filename, frame.Line)
		},
	}
	if err := os.Mkdir("../logs", 0644); err != nil {
		fmt.Printf("err creating log dir:%v", err)
		return
	}
	allFile, err := os.OpenFile("../logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		fmt.Printf("err creating log file:%v", err)
		return
	}
	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile,os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)

	return
}
