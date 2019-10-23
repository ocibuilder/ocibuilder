/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

// ParseBuildSpec parses the build specification which is read in through spec.yml
// or build.yaml and generates an array of build argumenets
func ParseBuildSpec(spec *v1alpha1.BuildSpec) ([]v1alpha1.ImageBuildArgs, error) {
	var imageBuilds []v1alpha1.ImageBuildArgs
	for _, step := range spec.Steps {

		step.Context = ValidateContext(step.Context)

		// TODO: support for more than just local context path in dockerfile generation
		dockerfilePath, err := GenerateDockerfile(step, spec.Templates, step.Context.LocalContext.ContextPath)

		// Perform cleanup of generated files if parse errors out
		if err != nil {
			for _, args := range imageBuilds {
				if err := os.Remove(args.Dockerfile); err != nil {
					log.WithError(err).Errorln("error cleaning up generated files")
				}
			}
			if err := os.Remove(dockerfilePath); err != nil {
				log.WithError(err).Errorln("error cleaning up generated files")
			}
			return nil, err
		}

		imageBuild := v1alpha1.ImageBuildArgs{
			Name:       step.Name,
			Tag:        step.Tag,
			Dockerfile: dockerfilePath,
			Purge:      step.Purge,
			Context:    step.Context,
		}
		imageBuilds = append(imageBuilds, imageBuild)

	}
	return imageBuilds, nil
}

// GenerateDockerfile takes in a build steps and generates a Dockerfile
// returns path to generated dockerfile
func GenerateDockerfile(step v1alpha1.BuildStep, templates []v1alpha1.BuildTemplate, contextPath string) (string, error) {
	var dockerfile []byte
	for idx, stage := range step.Stages {
		baseImage := parseBaseImage(stage.Base, stage.Name)

		if idx != 0 {
			baseImage = fmt.Sprintf("\n\n%s", baseImage)
		}
		dockerfile = append(dockerfile, baseImage...)

		// handles parsing of cmds in stage without a template
		tmp, err := parseCmdType(stage.Cmd)
		if err != nil {
			return "", err
		}
		dockerfile = append(dockerfile, tmp...)

		// handles parsing of cmds in a stage from a template
		for _, t := range templates {
			if stage.Template == t.Name {
				tmp, err := parseCmdType(t.Cmd)
				if err != nil {
					return "", err
				}
				dockerfile = append(dockerfile, tmp...)
			}
		}
	}

	// if path not specified set to current working directory
	if contextPath == "" {
		contextPath = "."
	}

	file, err := ioutil.TempFile(contextPath, "Dockerfile")
	if err != nil {
		return "", err
	}

	if _, err = file.Write(dockerfile); err != nil {
		return "", err
	}

	return filepath.Base(file.Name()), nil
}

// parseCmdType goes through a list of possible commands and parses them
// based on the request e.g. Docker/Ansible Path/Inline
func parseCmdType(cmds []v1alpha1.BuildTemplateStep) ([]byte, error) {
	var dockerfile []byte
	for _, cmd := range cmds {

		err := ValidateBuildTemplateStep(cmd)
		if err != nil {
			return nil, err
		}

		if cmd.Ansible != nil {
			tmp, err := ParseAnsibleCommands(cmd.Ansible)
			if err != nil {
				return nil, err
			}
			dockerfile = append(dockerfile, tmp...)
		}

		if cmd.Docker != nil {
			tmp, err := ParseDockerCommands(cmd.Docker)
			if err != nil {
				return nil, err
			}
			dockerfile = append(dockerfile, tmp...)
		}
	}
	return dockerfile, nil
}

// ParseAnsibleCommands is used to parse ansible commands from the ansible step
// and append the parsed template to Dockerfile
func ParseAnsibleCommands(ansibleStep *v1alpha1.AnsibleStep) ([]byte, error) {
	var buf bytes.Buffer
	var dockerfile []byte

	box := packr.NewBox("../templates/ansible")

	if ansibleStep.Local != nil {

		file, err := box.Find(v1alpha1.AnsiblePath)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("ansibleLocal").Parse(string(file))
		if err != nil {
			return nil, err
		}

		if err = tmpl.Execute(&buf, ansibleStep.Local); err != nil {
			return nil, err
		}
		dockerfileBytes := buf.Bytes()
		dockerfile = append(dockerfile, dockerfileBytes...)

		return dockerfile, nil
	}

	if ansibleStep.Galaxy != nil {
		file, err := box.Find(v1alpha1.AnsibleGalaxyPath)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New("ansibleGalaxy").Parse(string(file))
		if err != nil {
			return nil, err
		}

		if err = tmpl.Execute(&buf, ansibleStep.Galaxy); err != nil {
			return nil, err
		}
		dockerfileBytes := buf.Bytes()
		dockerfile = append(dockerfile, dockerfileBytes...)

		return dockerfile, nil
	}

	return nil, errors.New("no ansible galaxy or local definitions found")
}

// ParseDockerCommands parses the inputted docker commands and adds to dockerfile
func ParseDockerCommands(dockerStep *v1alpha1.DockerStep) ([]byte, error) {
	var dockerfile []byte

	if dockerStep.Inline != nil {
		return append(dockerfile, strings.Join(dockerStep.Inline, "\n")...), nil
	}

	if dockerStep.Path != "" {
		file, err := os.Open(dockerStep.Path)
		if err != nil {
			return nil, err
		}

		defer func() {
			if r := recover(); r != nil {
				log.Warnln("panic recovered to execute final cleanup", r)
			}
			if err := file.Close(); err != nil {
				log.WithError(err).Errorln("error closing file")
			}
		}()

		res, err := parser.Parse(file)
		if err != nil {
			return nil, err
		}

		var commands []v1alpha1.Command

		for _, child := range res.AST.Children {
			cmd := v1alpha1.Command{
				Cmd:       child.Value,
				Original:  child.Original,
				StartLine: child.StartLine,
				Flags:     child.Flags,
			}

			if child.Next != nil && len(child.Next.Children) > 0 {
				cmd.SubCmd = child.Next.Children[0].Value
				child = child.Next.Children[0]
			}

			cmd.IsJSON = child.Attributes["json"]
			for n := child.Next; n != nil; n = n.Next {
				cmd.Value = append(cmd.Value, n.Value)
			}
			commands = append(commands, cmd)
		}
		return addCommandsToDockerfile(commands, dockerfile), nil
	}

	return nil, errors.New("no docker cmd path or inline docker commands defined")
}

// parseBaseImage parses the base image specification to include image, platform
// and as conditions
func parseBaseImage(base v1alpha1.Base, name string) string {
	baseImage := fmt.Sprintf("FROM %s", base.Image)

	if base.Platform != "" {
		baseImage = fmt.Sprintf("%s:%s", baseImage, base.Platform)
	}

	if base.Tag != "" {
		baseImage = fmt.Sprintf("%s:%s", baseImage, base.Tag)
	}

	if name != "" {
		baseImage = fmt.Sprintf("%s AS %s", baseImage, name)
	}

	return fmt.Sprintf("%s\n", baseImage)
}

// addCommandsToDockerfile is used to append commands to dockerfile
func addCommandsToDockerfile(commands []v1alpha1.Command, dockerfile []byte) []byte {
	for _, command := range commands {
		cmd := command.Cmd

		if cmd == "from" {
			cmd = "\n" + cmd
		}
		line := strings.ToUpper(cmd) + " " + strings.Join(command.Value, " ") + "\n"
		dockerfile = append(dockerfile, line...)
	}
	return dockerfile
}
