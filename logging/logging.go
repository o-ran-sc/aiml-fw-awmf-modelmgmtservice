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
	"os"
)

var (
	LOG_LEVEL     = "LOG_LEVEL"
	LOG_FILE_NAME = "LOG_FILE_NAME"
)

func init() {

	var logLevel string
	val, flg := os.LookupEnv(LOG_LEVEL)
	if flg {
		logLevel = val
	} else {
		logLevel = "DEBUG"
	}
	Load(logLevel, os.Getenv(LOG_FILE_NAME))
	INFO("Loggers loaded ..")
}

// Prefixes INFO for each log message
func INFO(msg string, logParams ...interface{}) {
	Logger.Info(msg, logParams...)
}

// Prefixes WARN for each log message
func WARN(msg string, logParams ...interface{}) {
	Logger.Info(msg, logParams...)
}

// Prefixes ERROR for each log message
func ERROR(msg string, logParams ...interface{}) {
	Logger.Info(msg, logParams...)
}
