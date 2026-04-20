package reporter

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"parser/internal/model"
	"path/filepath"
)

// Reporter create reports
type Reporter struct {
	temp           *template.Template
	reportFilename string
}

// NewReporter create new instance
func NewReporter(templatePath string, reportFilename string) (*Reporter, error) {

	name := filepath.Base(templatePath)

	tmpl, err := template.New(name).
		ParseFiles(templatePath)

	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	return &Reporter{tmpl, reportFilename}, nil
}

// Report create incidents report to file
func (r *Reporter) Report(incidents []model.Incident) error {
	var b bytes.Buffer
	r.temp.Execute(&b, incidents)
	err := os.WriteFile(r.reportFilename, b.Bytes(), 0777)
	if err != nil {
		return err
	}
	return nil
}
