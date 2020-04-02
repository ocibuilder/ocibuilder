package generate

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"github.com/gobuffalo/packr"
)

func GenerateSpecification(generator v1alpha1.SpecGenerator, dry bool) error {
	spec, err := generator.Generate()
	if err != nil {
		return err
	}

	if dry {
		if _, err := os.Stdout.Write(spec); err != nil {
			return err
		}
		return nil
	}

	if err := ioutil.WriteFile("beval.yaml", spec, 0644); err != nil {
		return err
	}
	return nil
}

func generate(templateName string, templateSpec interface{}) ([]byte, error) {
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

	if err = tmpl.Execute(&buf, templateSpec); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
