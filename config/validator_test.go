package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWhenSuccess(t *testing.T) {
	configDataValidator := NewConfigDataValidator()
	manager := configManager{
		App: AppConfigData{
			MMES_URL:      "test",
			LOG_FILE_NAME: "test",
		},
		DB: DBConfigData{
			MODEL_FILE_POSTFIX: "test",
			INFO_FILE_POSTFIX:  "test",
			S3_URL:             "test",
			S3_ACCESS_KEY:      "test",
			S3_SECRET_KEY:      "test",
			S3_REGION:          "test",
		},
	}

	err := configDataValidator.validate(&manager)
	assert.Nil(t, err)
}

func TestValidateWhenFailedValidate(t *testing.T) {
	configDataValidator := NewConfigDataValidator()
	manager := configManager{}

	err := configDataValidator.validate(&manager)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrInvalidConfigData)
}

func TestValidateWhenFailedDBS3URL(t *testing.T) {
	configDataValidator := NewConfigDataValidator()
	manager := configManager{
		App: AppConfigData{
			MMES_URL:      "test",
			LOG_FILE_NAME: "test",
		},
		DB: DBConfigData{
			MODEL_FILE_POSTFIX: "test",
			INFO_FILE_POSTFIX:  "test",
			S3_ACCESS_KEY:      "test",
			S3_SECRET_KEY:      "test",
			S3_REGION:          "test",
		},
	}

	err := configDataValidator.validate(&manager)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrInvalidConfigData)
	assert.Equal(t, "", manager.DB.S3_URL)
}
