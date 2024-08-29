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
		c.errs = append(c.errs, fmt.Errorf("log_file_name: %s", manager.App.LOG_FILE_NAME))
	}

	if manager.App.MMES_URL == "" {
		c.errs = append(c.errs, fmt.Errorf("mmes_url: %s", manager.App.MMES_URL))
	}

	if manager.DB.MODEL_FILE_POSTFIX == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.MODEL_FILE_POSTFIX))
	}

	if manager.DB.INFO_FILE_POSTFIX == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.INFO_FILE_POSTFIX))
	}

	if manager.DB.S3_URL == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.S3_URL))
	}

	if manager.DB.S3_ACCESS_KEY == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.S3_ACCESS_KEY))
	}

	if manager.DB.S3_SECRET_KEY == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.S3_SECRET_KEY))
	}

	if manager.DB.S3_REGION == "" {
		c.errs = append(c.errs, fmt.Errorf("model_file_postfix: %s", manager.DB.S3_REGION))
	}

	return c.result()
}
