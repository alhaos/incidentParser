package interpreter

import (
	"log/slog"
	"parser/internal/model"
	"regexp"
)

type Interpreter struct {
	excludeRuleList *listNode
}

type listNode struct {
	ExcludeRule
	next *listNode
}

// Incident exclude rule
type ExcludeRule struct {
	Name     string
	Callback func(model.Incident) (bool, error)
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		excludeRuleList: nil,
	}
}

func (i *Interpreter) ShouldExclude(inc model.Incident) bool {

	if i.excludeRuleList == nil {
		return false
	}

	for currentNode := i.excludeRuleList; currentNode != nil; currentNode = currentNode.next {

		exclude, err := currentNode.ExcludeRule.Callback(inc)
		if err != nil {
			slog.Error("Interpreter rule", "error", err.Error())
			return false
		}

		if exclude {
			slog.Debug("Incident excluded", "id", inc.IncidentID, "rule", currentNode.ExcludeRule.Name)
			return true
		}
	}

	return false
}

func (i *Interpreter) AddRule(eRule ExcludeRule) {
	if i.excludeRuleList == nil {
		i.excludeRuleList = &listNode{
			ExcludeRule: eRule,
			next:        nil,
		}
		return
	}

	i.excludeRuleList = &listNode{
		ExcludeRule: eRule,
		next:        i.excludeRuleList,
	}
}

func (i *Interpreter) AddMessageMatchRule(pattern string, name string) {
	i.AddRule(
		ExcludeRule{
			Name: name,
			Callback: func(inc model.Incident) (bool, error) {
				match, err := regexp.MatchString(pattern, inc.Message)
				if err != nil {
					return false, err
				}

				if match {
					return true, nil
				}
				return false, nil
			},
		},
	)
}

func NewDefault(ruleSet model.RuleSet) *Interpreter {

	i := Interpreter{}

	i.AddRule(
		ExcludeRule{
			Name: "cleared",
			Callback: func(inc model.Incident) (bool, error) {
				if inc.Severity == "Clear" {
					return true, nil
				}
				return false, nil
			},
		},
	)

	for _, r := range ruleSet.MessageMatchRules {
		i.AddMessageMatchRule(r.Pattern, r.Name)
	}

	return &i
}
