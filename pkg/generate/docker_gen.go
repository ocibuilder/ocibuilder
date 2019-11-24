package generate

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/sirupsen/logrus"
)

type DockerGenerator struct {
	Filepath  string
	ImageName string
	Tag       string
	Logger    *logrus.Logger
}

func (d DockerGenerator) Generate() ([]byte, error) {
	file, err := os.Open(d.Filepath)
	if err != nil {
		return nil, err
	}

	res, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	var stages []string
	var templateName string
	tmpCmds := make(map[string][]string)

	templateIdx := 0
	for _, child := range res.AST.Children {

		cmd := v1alpha1.Command{
			Cmd:       child.Value,
			Original:  child.Original,
			StartLine: child.StartLine,
			Flags:     child.Flags,
		}

		for n := child.Next; n != nil; n = n.Next {
			cmd.Value = append(cmd.Value, n.Value)
		}

		if cmd.Cmd == "from" {
			templateIdx++
			templateName = "build-template-" + strconv.Itoa(templateIdx)

			fromCmd, _ := parseFromCmd(templateName, cmd)
			stage, err := generate("stage_tmpl", fromCmd)
			if err != nil {
				return nil, err
			}
			stages = append(stages, string(stage))

			tmpCmds[templateName] = make([]string, 0)
			continue
		}
		tmpCmds[templateName] = append(tmpCmds[templateName], cmd.Original)
	}

	var templates []string
	for templateName, templateCmds := range tmpCmds {
		buildTemplate := v1alpha1.BuildGenTemplate{Name: templateName, Cmds: templateCmds}
		tmpl, err := generate("build_template_tmpl", buildTemplate)
		if err != nil {
			return nil, err
		}
		templates = append(templates, string(tmpl))
	}

	ocibuilderTmpl := v1alpha1.GenerateTemplate{
		ImageName: d.ImageName,
		Stages:    stages,
		Tag:       d.Tag,
		Templates: templates,
	}
	spec, err := generate("ocibuilder_tmpl", ocibuilderTmpl)
	if err != nil {
		return nil, err
	}
	return removeWhitespace(spec), nil
}

func parseFromCmd(templateName string, cmd v1alpha1.Command) (v1alpha1.StageGenTemplate, error) {
	stageName := "build-stage"

	if len(cmd.Value) == 0 {
		return v1alpha1.StageGenTemplate{}, errors.New("malformed FROM command")
	}
	fields := cmd.Value
	image := strings.Split(fields[0], ":")

	if len(fields) > 1 && fields[1] == "AS" {
		stageName = fields[2]
	}

	stageTmp := v1alpha1.StageGenTemplate{
		StageName:    stageName,
		TemplateName: templateName,
		Base:         image[0],
	}

	if len(image) > 1 {
		stageTmp.BaseTag = image[1]
	}

	return stageTmp, nil
}

func removeWhitespace(spec []byte) []byte {
	re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
	specNoSpace := strings.Trim(re.ReplaceAllString(string(spec), ""), "\r\n")
	return []byte(specNoSpace)
}
