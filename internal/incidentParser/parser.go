package incidentParser

import (
	"os"
	"parser/internal/model"
	"parser/internal/utils"
	"strconv"
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(file string) (model.Incident, error) {

	incident := model.Incident{
		File: file,
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return model.Incident{}, err
	}

	for row := range strings.SplitSeq(string(data), "\r\n") {
		parts := strings.SplitN(row, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Host":
			incident.Host = value
		case "Target type":
			incident.TargetType = value
		case "Target name":
			incident.TargetName = strings.SplitN(value, "<", 2)[0]
		case "Incident creation time":
			incident.IncidentCreationTime, err = utils.ParseDataType1(value)
			if err != nil {
				return incident, err
			}
		case "Message":
			incident.Message = value
		case "Severity":
			incident.Severity = value
		case "Incident ID":
			incident.IncidentID, err = strconv.Atoi(value)
			if err != nil {
				return incident, err
			}
		}
	}
	return incident, nil
}
