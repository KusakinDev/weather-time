package loggerconfig

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {

	file, err := os.OpenFile("Logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	//defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	log.SetOutput(multiWriter)
}
