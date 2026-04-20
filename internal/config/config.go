package config

import (
	"os"

	"go.yaml.in/yaml/v4"
)

// Config general app config
type Config struct {
	TemplatePath   string `yaml:"templatePath"`
	ReportFilename string `yaml:"reportFilename"`
}

// NewConfig create new config
func NewConfig(filename string) (*Config, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Config

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
