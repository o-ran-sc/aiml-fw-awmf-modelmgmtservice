package config

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConfigData = errors.New("invalid config data")
)

func NewConfigDataValidator() *configDataValidator {
	return &configDataValidator{
		errs: []error{},
	}
}

type configDataValidator struct {
	errs []error
}

func (c *configDataValidator) result() error {
	var err error = nil
	if len(c.errs) > 0 {
		err = errors.Join(ErrInvalidConfigData, errors.Join(c.errs...))
	}

	return err
}

func (c *configDataValidator) validate(manager *configManager) error {
	if manager.App.LOG_FILE_NAME == "" {
		c.errs = append(c.errs, fmt.Errorf("log_file_name is not set/available or empty"))
	}

	if manager.App.MMES_URL == "" {
		c.errs = append(c.errs, fmt.Errorf("mmes_url is not set/available or empty"))
	}

	if manager.DB.MODEL_FILE_POSTFIX == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix is not set/available or empty"))
	}

	if manager.DB.INFO_FILE_POSTFIX == "" {
		c.errs = append(c.errs, fmt.Errorf("model_info_file_postfix is not set/available or empty"))
	}

	if manager.DB.S3_URL == "" {
		c.errs = append(c.errs, fmt.Errorf("s3_url is not set/available or empty"))
	}

	if manager.DB.S3_ACCESS_KEY == "" {
		c.errs = append(c.errs, fmt.Errorf("s3_access_key is not set/available or empty"))
	}

	if manager.DB.S3_SECRET_KEY == "" {
		c.errs = append(c.errs, fmt.Errorf("s3_secret_key is not set/available or empty"))
	}

	if manager.DB.S3_REGION == "" {
		c.errs = append(c.errs, fmt.Errorf("s3_region is not set/available or empty"))
	}

	return c.result()
}
