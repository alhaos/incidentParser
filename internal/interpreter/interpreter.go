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

func NewDefault() *Interpreter {

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

	rules := []model.Rule{
		{Pattern: "Total db time", Name: "blocking session"},
		{Pattern: "Tablespace.*is.*full", Name: "Tablespace is full"},
		{Pattern: "Users Lock Sessions", Name: "Users Lock Sessions"},
		{Pattern: "Agent Unreachable", Name: "Agent Unreachable"},
		{Pattern: "DBA_2", Name: "DBA_2"},
		{Pattern: "Hang replication session", Name: "Hang replication session"},
		{Pattern: "Stats Stale", Name: "Stats Stale"},
		{Pattern: "CDS Loader", Name: "CDS Loader"},
		{Pattern: "Snapshot not refresh", Name: "Snapshot not refresh"},
		{Pattern: "Hang jobs", Name: "Hang jobs"},
		{Pattern: "Capture need archivelog older than", Name: "Capture need archivelog older than"},
		{Pattern: "TM4_NO_PARSE_DATA_AT_LAST_ONE_SERVER", Name: "TM4_NO_PARSE_DATA_AT_LAST_ONE_SERVER"},
		{Pattern: "The value of DDL Locks is", Name: "The value of DDL Locks is"},
		{Pattern: "Agent is unable to communicate", Name: "Agent is unable to communicate"},
		{Pattern: "Invalid Object Count in", Name: "Invalid Object Count in"},
		{Pattern: "PREDIX Loader", Name: "PREDIX Loader"},
		{Pattern: "The value of Apply errors is", Name: "The value of Apply errors is"},
		{Pattern: "STREAMS error queue for apply process", Name: "STREAMS error queue for apply process"},
		{Pattern: "Standby database NOT_APPLIED logs", Name: "Standby database NOT_APPLIED logs"},
	}

	for _, r := range rules {
		i.AddMessageMatchRule(r.Pattern, r.Name)
	}

	return &i
}
