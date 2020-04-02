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

package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/beval/beval/ocictl/pkg/utils"
	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"github.com/beval/beval/pkg/common"
	"github.com/beval/beval/pkg/context"
	"github.com/beval/beval/pkg/request"
	"github.com/beval/beval/pkg/util"
	"github.com/beval/beval/pkg/validate"
	"github.com/gobuffalo/packr"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/pkg/errors"
)

// ParseBuildSpec parses the build specification which is read in through spec.yml
// or build.yaml and generates an array of build arguments
func ParseBuildSpec(spec *v1alpha1.BuildSpec) ([]v1alpha1.ImageBuildArgs, error) {
	var imageBuilds []v1alpha1.ImageBuildArgs
	kubeConfig, ok := os.LookupEnv(common.EnvVarKubeConfig)
	if !ok {
		kubeConfig = ""
	}
	for _, step := range spec.Steps {

		if err := validate.ValidateContext(step.BuildContext); err != nil {
			return nil, err
		}

		buildContext, err := context.GetBuildContextReader(step.BuildContext, kubeConfig)
		if err != nil {
			return nil, err
		}
		buildContextPath, err := buildContext.Read()
		if err != nil {
			return nil, err
		}
		cleanOnKill(buildContextPath)

		dockerfilePath, err := GenerateDockerfile(step, spec.Templates, buildContextPath+common.ContextDirectory)
		// Perform cleanup of generated files if parse errors out
		if err != nil {
			for _, args := range imageBuilds {
				if err := os.Remove(args.Dockerfile); err != nil {
					util.Logger.WithError(err).Errorln("error cleaning up generated files")
				}
			}
			if err := os.Remove(dockerfilePath); err != nil {
				util.Logger.WithError(err).Errorln("error cleaning up generated files")
			}
			return nil, err
		}

		if err := context.InjectDockerfile(buildContextPath, dockerfilePath); err != nil {
			return nil, errors.Errorf("error attempting to inject Dockerfile - err: %s", err)
		}

		imageBuild := v1alpha1.ImageBuildArgs{
			Name:             step.Name,
			Tag:              step.Tag,
			Dockerfile:       filepath.Base(dockerfilePath),
			Purge:            step.Purge,
			BuildContextPath: buildContextPath,
			Labels:           step.Labels,
			Creator:          step.Creator,
			Source:           step.Source,
			Cache:            step.Cache,
			StorageDriver:    spec.StorageDriver,
		}
		imageBuilds = append(imageBuilds, imageBuild)
	}
	return imageBuilds, nil
}

// GenerateDockerfile takes in a build steps and generates a Dockerfile
// returns path to generated dockerfile
func GenerateDockerfile(step v1alpha1.BuildStep, templates []v1alpha1.BuildTemplate, destination string) (string, error) {
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

	file, err := ioutil.TempFile(destination, "Dockerfile")
	if err != nil {
		return "", err
	}
	if _, err = file.Write(dockerfile); err != nil {
		return "", errors.Errorf("error generating Dockerfile err: %s", err)
	}
	return file.Name(), nil
}

// parseCmdType goes through a list of possible commands and parses them
// based on the request e.g. Docker/Ansible Path/Inline
func parseCmdType(cmds []v1alpha1.BuildTemplateStep) ([]byte, error) {
	var dockerfile []byte
	for _, cmd := range cmds {
		if err := validate.ValidateBuildTemplateStep(cmd); err != nil {
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

			if cmd.Docker.Inline != nil {
				return append(dockerfile, strings.Join(cmd.Docker.Inline, "\n")...), nil
			}

			if cmd.Docker.Path != "" {
				tmp, err := ParseDockerCommands(cmd.Docker.Path)
				if err != nil {
					return nil, err
				}
				dockerfile = append(dockerfile, tmp...)
			}

			if cmd.Docker.Url != "" {
				if err := request.RequestRemote(cmd.Docker.Url, common.DockerStepPath, cmd.Docker.Auth); err != nil {
					return nil, err
				}

				tmp, err := ParseDockerCommands(common.DockerStepPath)
				if err != nil {
					return nil, err
				}
				dockerfile = append(dockerfile, tmp...)
			}

		}

	}
	return dockerfile, nil
}

// ParseAnsibleCommands is used to parse ansible commands from the ansible step
// and append the parsed template to Dockerfile
func ParseAnsibleCommands(ansibleStep *v1alpha1.AnsibleStep) ([]byte, error) {
	buf := &bytes.Buffer{}
	var dockerfile []byte

	ansibleTemplateFunc := template.FuncMap{
		"DirExists": utils.DirExists,
		"newLine":   func() string { return "\n" },
	}

	// add newline to buffer before appending ansible commands
	//buf.WriteString("\n")
	box := packr.NewBox(v1alpha1.AnsibleTemplateDir)
	file, err := box.Find(v1alpha1.AnsibleTemplate)
	if err != nil {
		return nil, err
	}

	ansibleTemplate, err := template.New("Ansible").Funcs(ansibleTemplateFunc).Parse(string(file))
	if err != nil {
		return nil, err
	}

	if err := validate.SetAnsibleDefaultIfNotPresent(ansibleStep); err != nil {
		return nil, err
	}

	ansibleStep.Workspace = fmt.Sprintf("%s/%s", v1alpha1.AnsibleBase, ansibleStep.Workspace)
	if err := ansibleTemplate.Execute(buf, ansibleStep); err != nil {
		return nil, err
	}
	dockerfileBytes := buf.Bytes()
	dockerfile = append(dockerfile, dockerfileBytes...)
	return dockerfile, nil

}

// ParseDockerCommands parses the inputted docker commands and adds to dockerfile
func ParseDockerCommands(dockerCmdFilepath string) ([]byte, error) {
	var dockerfile []byte
	cmdFile, err := os.Open(dockerCmdFilepath)

	defer func() {
		if r := recover(); r != nil {
			util.Logger.Warnln("panic recovered to execute final cleanup", r)
		}
		if err := cmdFile.Close(); err != nil {
			util.Logger.WithError(err).Errorln("error closing cmdFile")
		}
	}()

	if err != nil {
		return nil, err
	}

	res, err := parser.Parse(cmdFile)
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

func cleanOnKill(contextPath string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := os.RemoveAll(contextPath + "/ocib"); err != nil {
			fmt.Println("error cleaning up files", err)
		}
		os.Exit(1)
	}()
}
