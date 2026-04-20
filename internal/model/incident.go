package model

import "time"

type VerdictStatus int

const (
	undefined = iota
	important
	trash
)

type Incident struct {
	File                 string
	Host                 string
	TargetType           string
	TargetName           string
	IncidentCreationTime time.Time
	Message              string
	Severity             string
	IncidentID           int
	Verdict              VerdictStatus
}
