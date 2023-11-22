package logging

import (
	"log"
	"os"
)

var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger

func init() {
	logFile, fileErr := os.Create(os.Getenv("LOG_FILE_NAME"))
	if fileErr != nil {
		log.Fatal("Can not start MMES service,issue in creating log file")
	}
	flags := log.Ldate | log.Ltime 
	infoLogger = log.New(logFile, "INFO:", flags)
	warnLogger = log.New(logFile, "WARN:", flags)
	errorLogger = log.New(logFile, "ERROR:", flags)

	INFO("Loggers loaded ..")
}

// Prefixes INFO for each log message
func INFO(logParams ...interface{}) {
	infoLogger.Println(logParams...)
}

// Prefixes WARN for each log message
func WARN(logParams ...interface{}) {
	warnLogger.Println(logParams...)
}

// Prefixes ERROR for each log message
func ERROR(logParams ...interface{}) {
	errorLogger.Println(logParams...)
}
