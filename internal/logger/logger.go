package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

const logPath = "/app/logs/base.log"

func InitLogger(logLevel, mode string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Ошибка при парсинге уровня логирования: %v", err)
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)

	if mode == "prod" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return fmt.Sprintf("%s:%d ", frame.Function, frame.Line), ""
			},
		})

		if err := os.MkdirAll("/app/logs", 0777); err != nil {
			logrus.Fatal("Ошибка создания папки логов: ", err)
		}

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			logrus.Fatal("Ошибка открытия файла логов: ", err)
		}
		logrus.SetOutput(logFile)
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				return fmt.Sprintf("%s:%d ", frame.Function, frame.Line), ""
			},
		})
	}

}
