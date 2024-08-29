package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadWhenSuccess(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		ENV_KEY_DB_S3_SECRET_KEY:      "secret",
		ENV_KEY_DB_INFO_FILE_POSTFIX:  "_test.info",
		ENV_KEY_DB_MODEL_FILE_POSTFIX: "_test.zip",
	}
	unsetFunc := osSetData(defaultEnvData)
	defer unsetFunc()

	Load(nil, NewEnvDataLoader(nil))
	assert.Equal(t, "secret", manager.DB.S3_SECRET_KEY)
	assert.Equal(t, "_test.info", manager.DB.INFO_FILE_POSTFIX)
	assert.Equal(t, "_test.zip", manager.DB.MODEL_FILE_POSTFIX)
}

func TestLoadWhenCalledTwice(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		ENV_KEY_DB_S3_SECRET_KEY: "secret",
	}
	unsetFunc := osSetData(defaultEnvData)
	defer unsetFunc()

	Load(nil, NewEnvDataLoader(nil))

	assert.Equal(t, "secret", manager.DB.S3_SECRET_KEY)

	defaultEnvData = DefaultEnvData{
		ENV_KEY_DB_S3_SECRET_KEY: "secret2",
	}
	unsetFunc = osSetData(defaultEnvData)
	defer unsetFunc()
	Load(nil, NewEnvDataLoader(nil))

	assert.Equal(t, "secret", manager.DB.S3_SECRET_KEY)
}

func TestGetConfigManagerWhenTriedToOverWrite(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		ENV_KEY_DB_S3_SECRET_KEY: "secret",
	}
	unsetFunc := osSetData(defaultEnvData)
	defer unsetFunc()

	configManagerInstance := GetConfigManager()

	assert.Equal(t, "secret", configManagerInstance.DB.S3_SECRET_KEY)

	configManagerInstance.DB.S3_SECRET_KEY = "secret2"

	configManagerInstance = GetConfigManager()
	assert.Equal(t, "secret", configManagerInstance.DB.S3_SECRET_KEY)
}

// Please check plain text log manually
func TestConfigManagerString(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		ENV_KEY_DB_S3_SECRET_KEY: "secret",
	}
	unsetFunc := osSetData(defaultEnvData)
	defer unsetFunc()

	configManagerInstance := GetConfigManager()
	summary := configManagerInstance.String()

	assert.NotEmpty(t, summary)
}
