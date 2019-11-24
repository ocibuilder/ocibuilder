package generate

import (
	"bytes"
	"html/template"

	"github.com/gobuffalo/packr"
)

func Generate(templateName string, specTmpl interface{}) ([]byte, error) {
	var buf bytes.Buffer
	box := packr.NewBox("../../config/templates")

	file, err := box.Find(templateName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("generatedSpec").Parse(string(file))
	if err != nil {
		return nil, err
	}

	if err = tmpl.Execute(&buf, specTmpl); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
