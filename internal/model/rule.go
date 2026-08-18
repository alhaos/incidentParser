package model

type Rule struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
}

type MessageMatchRule struct {
	Rule
}

type RuleSet struct {
	MessageMatchRules []MessageMatchRule `yaml:"messageMatchRules"`
}
