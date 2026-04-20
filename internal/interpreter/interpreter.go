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

func Default() *Interpreter {

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

	rules := []struct {
		pattern string
		name    string
	}{
		{"Total db time", "blocking session"},
		{"Tablespace.*is.*full", "Tablespace is full"},
		{"Users Lock Sessions", "Users Lock Sessions"},
		{"Agent Unreachable", "Agent Unreachable"},
		{"DBA_2", "DBA_2"},
		{"Hang replication session", "Hang replication session"},
		{"Stats Stale", "Stats Stale"},
		{"CDS Loader", "CDS Loader"},
		{"Snapshot not refresh", "Snapshot not refresh"},
		{"Hang jobs", "Hang jobs"},
		{"Capture need archivelog older than", "Capture need archivelog older than"},
		{"TM4_NO_PARSE_DATA_AT_LAST_ONE_SERVER", "TM4_NO_PARSE_DATA_AT_LAST_ONE_SERVER"},
		{"The value of DDL Locks is", "The value of DDL Locks is"},
		{"Agent is unable to communicate", "Agent is unable to communicate"},
		{"Invalid Object Count in", `Invalid Object Count in`},
		{"PREDIX Loader", "PREDIX Loader"},
		{"The value of Apply errors is", "The value of Apply errors is"},
		{"STREAMS error queue for apply process", "STREAMS error queue for apply process"},
		{"Standby database NOT_APPLIED logs", "Standby database NOT_APPLIED logs"},
	}

	for _, r := range rules {
		i.AddMessageMatchRule(r.pattern, r.name)
	}

	return &i
}
