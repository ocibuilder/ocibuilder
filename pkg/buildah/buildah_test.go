/*
Copyright © 2019 BlackRock Inc.

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

package buildah

import (
	"os"
	"os/exec"
	"testing"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/fake"
	"github.com/stretchr/testify/assert"
)

// TestBuildah_Build is the test for a buildah build (or bud)
func TestBuildah_Build(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	b := Buildah{
		Logger: common.GetLogger(true),
	}
	_, err := b.Build(fake.Spec)
	assert.Equal(t, nil, err)
	b.Clean()
}

// TestBuildah_Login is the test for a buildah login
func TestBuildah_Login(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	b := Buildah{
		Logger: common.GetLogger(true),
	}
	_, err := b.Login(fake.Spec)
	assert.Equal(t, nil, err)
}

// TestBuildah_Push is the test for a buildah push
func TestBuildah_Push(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	b := Buildah{
		Logger: common.GetLogger(true),
	}
	_, err := b.Push(fake.Spec)
	assert.Equal(t, nil, err)
}

// TestBuildah_Pull is the test for a buildah pull
func TestBuildah_Pull(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	b := Buildah{
		Logger: common.GetLogger(true),
	}
	_, err := b.Pull(fake.Spec, "alpine")
	assert.Equal(t, nil, err)
}

func TestCreateBuildCommand(t *testing.T) {
	expectedBuildCommand := []string{"bud", "-f", "path/to/Dockerfile", "-t", "image-name:1.0.0", "./ocib/context/context.tar.gz"}

	buildCommand := createBuildCommand(buildArgs, "")
	assert.Equal(t, expectedBuildCommand, buildCommand)
}

func TestCreateBuildCommandStorageDriver(t *testing.T) {
	expectedBuildCommand := []string{"bud", "-f", "path/to/Dockerfile", "--storage-driver", "vfs", "-t", "image-name:1.0.0", "./ocib/context/context.tar.gz"}

	buildCommand := createBuildCommand(buildArgs, "vfs")
	assert.Equal(t, expectedBuildCommand, buildCommand)
}

func TestCreateLoginCommand(t *testing.T) {
	expectedLoginCommand := []string{"login", "-u", "username", "-p", "ThiSiSalOgInToK3N", "example-registry"}

	loginCommand, err := createLoginCommand(fake.LoginSpec[0])
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedLoginCommand, loginCommand)
}

func TestCreatePullCommand(t *testing.T) {
	expectedPullCommand := []string{"pull", "--creds", "username:ThiSiSalOgInToK3N", "example-registry/test-image:1.0.0"}
	b := Buildah{
		Logger: common.GetLogger(true),
	}
	pullCommand, err := b.createPullCommand("example-registry", "test-image:1.0.0", fake.Spec)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedPullCommand, pullCommand)
}

func TestCreatePushCommand(t *testing.T) {
	expectedPushCommand := []string{"push", "--creds", "username:ThiSiSalOgInToK3N", "example-registry/example-image:1.0.0"}
	b := Buildah{
		Logger: common.GetLogger(true),
	}
	pushCommand, err := b.createPushCommand("example-registry", "example-image:1.0.0", fake.Spec)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedPushCommand, pushCommand)
}

var buildArgs = v1alpha1.ImageBuildArgs{
	Name:             "image-name",
	Tag:              "1.0.0",
	Dockerfile:       "path/to/Dockerfile",
	BuildContextPath: ".",
}

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
