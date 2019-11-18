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

package docker

import (
	"testing"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/dummy"
	"github.com/stretchr/testify/assert"
)

var docker = Docker{
	Logger: common.GetLogger(true),
	Client: testClient{},
}

// TestDocker_Build is the test for a docker build
func TestDocker_Build(t *testing.T) {
	_, err := docker.Build(dummy.Spec)
	assert.Equal(t, nil, err)
	docker.Clean()
}

// TestDocker_Login is the test for a docker login
func TestDocker_Login(t *testing.T) {
	_, err := docker.Login(dummy.Spec)
	assert.Equal(t, nil, err)
}

// TestDocker_Push is the test for a docker push
func TestDocker_Push(t *testing.T) {
	_, err := docker.Push(dummy.Spec)
	assert.Equal(t, nil, err)
}

// TestDocker_Pull is the test for a docker pull
func TestDocker_Pull(t *testing.T) {
	_, err := docker.Pull(dummy.Spec, "testImage")
	assert.Equal(t, nil, err)
}
