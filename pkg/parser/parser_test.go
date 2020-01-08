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
	"io/ioutil"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/stretchr/testify/assert"
)

const expectedAnsibleLocalCommands = `
ENV PLAYBOOK_DIR /etc/ansible/
RUN mkdir -p $PLAYBOOK_DIR
WORKDIR $PLAYBOOK_DIR
COPY templates templates
COPY files files
COPY vars vars
COPY tasks tasks
ADD *.yaml ./
RUN ansible-playbook /playbook.yaml`

const expectedAnsibleGalaxyCommands = `
ENV PLAYBOOK_DIR /etc/ansible/
RUN mkdir -p $PLAYBOOK_DIR
WORKDIR $PLAYBOOK_DIR
COPY templates templates
COPY files files
COPY vars vars
COPY tasks tasks
ADD *.yaml ./
RUN if [ -f /requirements.yaml ]; then ansible-galaxy install -r /requirements.yaml; fi
RUN ansible-galaxy install TestGalaxy`

const expectedInlineDockerfile = "FROM go / java / nodejs / python:ubuntu_xenial:v1.0.0 AS first-stage\nADD ./ /test-path\nWORKDIR /test-dir\nENV PORT=3001\nCMD [\"go\", \"run\", \"main.go\"]\n\nFROM alpine:latest AS second-stage\nCMD [\"echo\", \"done\"]"

const expectedDockerfile = "FROM go / java / nodejs / python:ubuntu_xenial:v1.0.0 AS first-stage\nRUN pip install kubernetes\nCOPY app/ /bin/app\n\n\nFROM alpine:latest AS second-stage\nCMD [\"echo\", \"done\"]"

func TestParseDockerCommands(t *testing.T) {
	path := "../../testing/dummy/commands_basic_parser_test.txt"
	dockerfile, err := ParseDockerCommands(path)
	expectedDockerfile := "RUN pip install kubernetes\nCOPY app/ /bin/app\n"

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedDockerfile, string(dockerfile))
}

func TestGenerateDockerfile(t *testing.T) {
	file, err := ioutil.ReadFile("../../testing/dummy/build.yaml")
	assert.Equal(t, nil, err)

	templates := []v1alpha1.BuildTemplate{{
		Name: "template-1",
		Cmd: []v1alpha1.BuildTemplateStep{{
			Docker: &v1alpha1.DockerStep{
				Path: "../../testing/dummy/commands_basic_parser_test.txt",
			},
		},
		}}}

	buildSpecification := v1alpha1.BuildSpec{}
	if err := yaml.Unmarshal(file, &buildSpecification); err != nil {
		assert.Fail(t, "fail unmarshalling build spec")
	}

	path, err := GenerateDockerfile(buildSpecification.Steps[0], templates, "")
	assert.Equal(t, nil, err)

	defer os.Remove(path)

	dockerfile, err := ioutil.ReadFile(path)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedDockerfile, string(dockerfile))
}

func TestGenerateDockerfileInline(t *testing.T) {
	file, err := ioutil.ReadFile("../../testing/dummy/build.yaml")
	assert.Equal(t, nil, err)

	template := []v1alpha1.BuildTemplate{{
		Name: "template-1",
		Cmd: []v1alpha1.BuildTemplateStep{{
			Docker: &v1alpha1.DockerStep{
				Inline: []string{
					"ADD ./ /test-path",
					"WORKDIR /test-dir",
					"ENV PORT=3001",
					"CMD [\"go\", \"run\", \"main.go\"]",
				},
			},
		}},
	}}

	buildSpecification := v1alpha1.BuildSpec{}
	if err := yaml.Unmarshal(file, &buildSpecification); err != nil {
		assert.Fail(t, "fail unmarshalling build spec")
	}

	path, err := GenerateDockerfile(buildSpecification.Steps[0], template, "")
	assert.Equal(t, nil, err)
	defer os.Remove(path)

	dockerfile, err := ioutil.ReadFile(path)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedInlineDockerfile, string(dockerfile))

	path, err = GenerateDockerfile(buildSpecification.Steps[0], template, "../../testing/dummy")
	assert.Equal(t, nil, err)

	dockerfile, err = ioutil.ReadFile("../../testing/dummy/" + path)
	assert.Equal(t, nil, err)
	defer os.Remove("../../testing/dummy/" + path)

	assert.Equal(t, expectedInlineDockerfile, string(dockerfile))
}

func TestParseAnsibleCommands(t *testing.T) {
	ansibleStep := &v1alpha1.AnsibleStep{
		Local: &v1alpha1.AnsibleLocal{
			Playbook: "/playbook.yaml",
		},
	}
	dockerfile, err := ParseAnsibleCommands(ansibleStep)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedAnsibleLocalCommands, string(dockerfile), "The generated ansible local commands must match expected")
}

func TestParseAnsibleGalaxyCommands(t *testing.T) {
	ansibleStepGalaxy := &v1alpha1.AnsibleStep{
		Galaxy: &v1alpha1.AnsibleGalaxy{
			Name:         "TestGalaxy",
			Requirements: "/requirements.yaml",
		},
	}
	dockerfile, err := ParseAnsibleCommands(ansibleStepGalaxy)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedAnsibleGalaxyCommands, string(dockerfile), "The generated ansible galaxy comnmands must match expected")
}
