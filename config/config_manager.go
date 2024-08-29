package config

import (
	"bytes"
	"sync"
)

var (
	manager = &configManager{}
	once    sync.Once
)

func GetConfigManager() configManager {
	return *manager
}

// This function is executed only once, after execution you can't load config data
func Load(loaders ...configLoader) {
	once.Do(func() {
		for _, loader := range loaders {
			loader.load(manager)
		}
	})
}

type configLoader interface {
	load(*configManager)
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
