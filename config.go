package coreconfig

import (
	"encoding/json"
	"io"
	"sync"
)

type CoreConfig struct {
	services []Service
	config   interface{}
}

type Service interface {
	Initialize(config interface{}) error
}

var (
	configInstance *CoreConfig
	once           sync.Once
)

func GetConfigInstance() *CoreConfig {
	once.Do(func() {
		configInstance = &CoreConfig{}
	})
	return configInstance
}

func (c *CoreConfig) Register(service Service) {
	c.services = append(c.services, service)
}

func (c *CoreConfig) LoadConfig(reader io.Reader, configTemplate interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(configTemplate); err != nil {
		return err
	}
	c.config = configTemplate
	return nil
}

func (c *CoreConfig) InitializeServices() error {
	for _, service := range c.services {
		if err := service.Initialize(c.config); err != nil {
			return err
		}
	}
	return nil
}
