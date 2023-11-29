/*
==================================================================================
Copyright (c) 2023 Samsung Electronics Co., Ltd. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
==================================================================================
*/
package logging

import (
	"log"
	"os"
)

var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger

func init() {
	//TODO add current timestamp as file prefix to retain the old log file
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
