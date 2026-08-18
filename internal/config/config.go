package config

import (
	"os"
	"parser/internal/model"

	"go.yaml.in/yaml/v4"
)

// Config general app config
type Config struct {
	TemplatePath   string        `yaml:"templatePath"`
	ReportFilename string        `yaml:"reportFilename"`
	IncidentsPath  string        `yaml:"incidentsPath"`
	RuleSet        model.RuleSet `yaml:"ruleSet"`
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
