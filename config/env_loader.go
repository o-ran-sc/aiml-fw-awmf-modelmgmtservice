package config

import (
	"fmt"
	"os"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
	"github.com/joho/godotenv"
)

// DB ENV KEY
const (
	ENV_KEY_DB_S3_URL             = "S3_URL"
	ENV_KEY_DB_S3_ACCESS_KEY      = "S3_ACCESS_KEY"
	ENV_KEY_DB_S3_SECRET_KEY      = "S3_SECRET_KEY"
	ENV_KEY_DB_S3_REGION          = "S3_REGION"
	ENV_KEY_DB_INFO_FILE_POSTFIX  = "INFO_FILE_POSTFIX"
	ENV_KEY_DB_MODEL_FILE_POSTFIX = "MODEL_FILE_POSTFIX"
)

// APP ENV KEY
const (
	ENV_KEY_APP_MMES_URL      = "MMES_URL"
	ENV_KEY_APP_LOG_FILE_NAME = "LOG_FILE_NAME"
)

type DefaultEnvData map[string]string

type envDataLoader struct {
}

func NewEnvDataLoader(defaultData DefaultEnvData, envFilePath ...string) *envDataLoader {
	if len(envFilePath) != 0 {
		err := godotenv.Load(envFilePath...)
		logging.ERROR(fmt.Sprintf("Failed to load env file, error msg: %s", err))
	}

	for k, v := range defaultData {
		if os.Getenv(k) != "" {
			continue
		}

		os.Setenv(k, v)
	}

	return &envDataLoader{}
}

func (e *envDataLoader) load(c *configManager) {
	e.appDataLoad(c)
	e.dbDataLoad(c)
}

func (e *envDataLoader) dbDataLoad(c *configManager) {
	c.DB.S3_URL = getStringWithDefault(ENV_KEY_DB_S3_URL, "http://localhost:32080")
	c.DB.S3_ACCESS_KEY = getStringWithDefault(ENV_KEY_DB_S3_ACCESS_KEY, "")
	c.DB.S3_SECRET_KEY = getStringWithDefault(ENV_KEY_DB_S3_SECRET_KEY, "")
	c.DB.S3_REGION = getStringWithDefault(ENV_KEY_DB_S3_REGION, "")
	c.DB.INFO_FILE_POSTFIX = getStringWithDefault(ENV_KEY_DB_INFO_FILE_POSTFIX, "_info.json")
	c.DB.MODEL_FILE_POSTFIX = getStringWithDefault(ENV_KEY_DB_MODEL_FILE_POSTFIX, "_model.zip")
}

func (e *envDataLoader) appDataLoad(c *configManager) {
	c.App.MMES_URL = getStringWithDefault(ENV_KEY_APP_MMES_URL, "0.0.0.0:8083")
	c.App.LOG_FILE_NAME = getStringWithDefault(ENV_KEY_APP_LOG_FILE_NAME, "mmes.log")
}

func getStringWithDefault(key string, defaultValue string) string {
	data := os.Getenv(key)
	if data == "" {
		return defaultValue
	}

	return data
}
