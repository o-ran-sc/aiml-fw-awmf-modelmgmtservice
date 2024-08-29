package config

import (
	"bytes"
	"sync"
)

var (
	manager = &configManager{}
	once    sync.Once
)

type configValidator interface {
	validate(*configManager) error
}

type configLoader interface {
	load(*configManager)
}

func GetConfigManager() configManager {
	return *manager
}

// This function is executed only once, after execution you can't load config data
func Load(validator configValidator, loaders ...configLoader) error {
	var err error = nil
	once.Do(func() {
		for _, loader := range loaders {
			loader.load(manager)
		}
		if validator != nil {
			err = validator.validate(manager)
		}
	})
	return err
}

type configManager struct {
	App AppConfigData
	DB  DBConfigData
}

func (c configManager) String() string {
	var buf bytes.Buffer
	buf.WriteString("------------Config Data------------")
	buf.WriteString("\n")
	buf.WriteString(c.App.String())
	buf.WriteString("\n")
	buf.WriteString(c.DB.String())
	buf.WriteString("\n")
	buf.WriteString("-----------------------------------")
	return buf.String()
}
