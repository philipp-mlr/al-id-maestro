package model

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	file, err := os.ReadFile("./config.yml")
	if err != nil {
		return fmt.Errorf("error opening config.yml: %v", err)
	}

	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return fmt.Errorf("error reading config.yml: %v", err)
	}

	log.Println("Loaded configuration")

	return nil
}
