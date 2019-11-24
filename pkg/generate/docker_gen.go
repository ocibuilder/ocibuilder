package generate

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

type DockerGenerator struct {
	Filepath string
}

func (d DockerGenerator) Generate() (interface{}, error) {
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
			stage, err := Generate("stage_tmpl", fromCmd)
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
		tmpl, err := Generate("build_template_tmpl", buildTemplate)
		if err != nil {
			return nil, err
		}
		templates = append(templates, string(tmpl))
	}

	ocibuilderTmpl := v1alpha1.GenerateTemplate{
		Stages:    stages,
		Templates: templates,
	}
	byt, _ := Generate("ocibuilder_tmpl", ocibuilderTmpl)
	fmt.Println(string(byt))

	return nil, nil
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
