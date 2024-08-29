package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func createTempEnvFile(filename string, envDataMap map[string]string) (func(), error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	for k, v := range envDataMap {
		file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return func() {
		os.Remove(filename)
	}, nil
}

func osSetData(envDataMap map[string]string) func() {
	for k, v := range envDataMap {
		os.Setenv(k, v)
	}

	return func() {
		for k := range envDataMap {
			os.Unsetenv(k)
		}
	}
}

func unSetViperDefault(envDataMap map[string]string) {
	for k := range envDataMap {
		viper.SetDefault(k, nil)
	}
}

func TestNewEnvDataLoadWhenGivenDefaultEnvFilePath(t *testing.T) {
	envFilePath := "./test.env"

	envDataMap := map[string]string{
		"testkey1": "test-value1",
		"testkey2": "test-value2",
	}
	deleteFunc, err := createTempEnvFile(envFilePath, envDataMap)
	if err != nil {
		t.Errorf("Failed to generate envfile %s\n", err)
		return
	}
	defer deleteFunc()

	NewEnvDataLoader(nil, envFilePath)
	assert.Equal(t, envDataMap["testkey1"], viper.GetString("testkey1"))
	assert.Equal(t, envDataMap["testkey2"], viper.GetString("testkey2"))
	assert.NotEqual(t, envDataMap["testkey1"], viper.GetString("testkey2"))
	assert.Empty(t, envDataMap["testkey3"])
}

func TestNewEnvDataLoadWhenGivenDefaultEnvData(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		"testkey1":        "test-value1",
		ENV_KEY_DB_S3_URL: "google",
	}
	NewEnvDataLoader(defaultEnvData)
	defer unSetViperDefault(defaultEnvData)

	assert.Equal(t, defaultEnvData["testkey1"], viper.GetString("testkey1"))
	assert.Equal(t, defaultEnvData[ENV_KEY_DB_S3_URL], viper.GetString(ENV_KEY_DB_S3_URL))
	assert.NotEqual(t, defaultEnvData["testkey1"], viper.GetString(ENV_KEY_DB_S3_URL))
	assert.Empty(t, defaultEnvData["testkey3"])
}

func TestNewEnvDataLoadWhenSuccess(t *testing.T) {
	os.Setenv(ENV_KEY_DB_S3_URL, "google")
	os.Setenv(ENV_KEY_APP_LOG_FILE_NAME, "test.log")

	unsetFunc := osSetData(map[string]string{
		ENV_KEY_DB_S3_URL:         "google",
		ENV_KEY_APP_LOG_FILE_NAME: "test.log",
	})
	defer unsetFunc()

	loader := NewEnvDataLoader(nil)
	configManager := configManager{}
	loader.load(&configManager)

	assert.Equal(t, "google", configManager.DB.S3_URL)
	assert.Equal(t, "test.log", configManager.App.LOG_FILE_NAME)

}

func TestNewEnvDataLoaderWhenFailedReadFile(t *testing.T) {
	defaultEnvData := DefaultEnvData{
		"testkey": "test-value1",
	}
	NewEnvDataLoader(defaultEnvData, "test.env")
	defer unSetViperDefault(defaultEnvData)

	assert.Equal(t, "test-value1", viper.GetString("testkey"))
}

func TestGetStringWithDefaultWhenSuccess(t *testing.T) {
	data := getStringWithDefault("testkey", "default")

	assert.Equal(t, data, "default")
}
