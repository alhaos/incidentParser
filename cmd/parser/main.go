package main

import (
	"flag"
	"log/slog"
	"os"
	"parser/internal/config"
	"parser/internal/incidentParser"
	"parser/internal/interpreter"
	"parser/internal/model"
	"parser/internal/reporter"
	"path/filepath"
	"strings"
)

const incidentsPath = `data/Incidents`
const templatePath = `template/template.html`
const defaultConfigPath = `config/config.yml`
const fileExtension = ".txt"

func main() {

	// Parse config filename flag
	configFilenamePointer := flag.String("config", defaultConfigPath, "-config config/config.yml")
	flag.Parse()
	configFilename := *configFilenamePointer

	// Init configuration
	conf, err := config.NewConfig(configFilename)
	if err != nil {
		panic(err)
	}

	// init log
	l := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			nil,
		),
	)

	slog.SetDefault(l)

	// Get data files list
	files, err := os.ReadDir(incidentsPath)

	if err != nil {
		panic(err)
	}

	// Init parser
	p := incidentParser.NewParser()

	// Init interpreter
	i := interpreter.Default()

	var importantIncidents []model.Incident
	var counter int

	// Loop files
	for _, entry := range files {

		name := entry.Name()

		if entry.IsDir() {
			slog.Debug("found: ", "name", name, "skip reason", "is dir")
			continue
		}

		counter++

		ext := filepath.Ext(name)
		lowerExt := strings.ToLower(ext)
		if lowerExt != strings.ToLower(fileExtension) {
			slog.Debug("found: ", "name", name, "skip reason", "Incorrect file extension")
			continue
		}

		incident, err := p.Parse(filepath.Join(incidentsPath, entry.Name()))
		if err != nil {
			panic(err)

		}

		if i.ShouldExclude(incident) {
			slog.Debug("found: ", "name", name, "skip reason", "unimportant")
			continue
		}

		slog.Info("found: ", "name", name)
		importantIncidents = append(importantIncidents, incident)
	}

	slog.Info("data files found", "count", counter)
	slog.Info("important incidents found", "count", len(importantIncidents))

	// Init reporter
	r, err := reporter.NewReporter(conf.TemplatePath, conf.ReportFilename)
	if err != nil {
		panic(err)
	}

	// Make report
	r.Report(importantIncidents)

}
