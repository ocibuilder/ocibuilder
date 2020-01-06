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

package command

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var builder = Builder("buildah")

func TestCommand_constructCommand(t *testing.T) {
	constructedCmd := cmd.ConstructCommand()

	assert.Equal(t, []string{"build", "-f", "./Dockerfile", "--storage-driver", "vfs", ".", "one", "two"}, constructedCmd)
}

func TestCommand_Exec(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	_, _, err := cmd.Exec()
	assert.Equal(t, nil, err)
}

func TestCommandBuilder_Flags(t *testing.T) {

	flags := []Flag{
		{"f", "Dockerfile", true, true},
		{"storage-driver", "", false, true},
		{"t", "image:tag", true, true},
	}

	builder := Builder("test").SetFlags(flags...)
	assert.Equal(t, Command{
		Name:    "test",
		Command: "",
		Flags:   expectedFlags,
		Args:    []string{},
	}, builder.Build())
}

var expectedFlags = []Flag{
	{"f", "Dockerfile", true, true},
	{"t", "image:tag", true, true},
}

var cmd = builder.SetCommand("build").SetFlags([]Flag{
	{
		Name:  "f",
		Value: "./Dockerfile",
		Short: true,
	},
	{
		Name:  "storage-driver",
		Value: "vfs",
		Short: false,
	},
}...).SetArgs(".", "one", "two").Build()

// enabling the mocking of exec commands as in https://npf.io/2015/06/testing-exec-command/
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
