package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"parser/internal/config"
	"parser/internal/incidentParser"
	"parser/internal/interpreter"
	"parser/internal/model"
	"parser/internal/reporter"
	"path/filepath"
)

const incidentsPath = `data/Incidents`
const templatePath = `template/template.html`
const defaultConfigPath = `config/config.yml`

func main() {

	configFilenamePointer := flag.String("config", defaultConfigPath, "-config config/config.yml")
	flag.Parse()
	configFilename := *configFilenamePointer

	fmt.Println(os.Getwd())

	conf, err := config.NewConfig(configFilename)
	if err != nil {
		panic(err)
	}

	l := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)

	slog.SetDefault(l)

	files, err := os.ReadDir(incidentsPath)

	if err != nil {
		panic(err)
	}

	p := incidentParser.NewParser()

	i := interpreter.Default()

	var importantIncidents []model.Incident

	for _, file := range files {

		incident, err := p.Parse(filepath.Join(incidentsPath, file.Name()))
		if err != nil {
			panic(err)

		}

		if !i.ShouldExclude(incident) {
			importantIncidents = append(importantIncidents, incident)
		}
	}

	r, err := reporter.NewReporter(conf.TemplatePath, conf.ReportFilename)
	if err != nil {
		panic(err)
	}

	r.Report(importantIncidents)

}

func printIncident(inc model.Incident) {
	fmt.Println("Incident:")
	fmt.Println("  ID:", inc.IncidentID)
	fmt.Println("  Type:", inc.TargetType)
	fmt.Println("  Name:", inc.TargetName)
	fmt.Println("  Host:", inc.Host)
	fmt.Println("  Time:", inc.IncidentCreationTime)
	fmt.Println("  Message:", inc.Message)
}
