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

package oci

import (
	"os"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/types"
	"github.com/ocibuilder/ocibuilder/testing/dummy"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	builder := Builder{
		Logger:   common.GetLogger(true),
		Client:   testClient{},
		Metadata: []v1alpha1.ImageMetadata{},
	}

	res := make(chan types.OCIBuildResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
		close(finished)
	}()

	go builder.Build(&dummy.Spec, res, errChan, finished)

	for {
		select {
		case err := <-errChan:
			{
				assert.Equal(t, nil, err)
				return
			}
		case <-res:
			{
			}
		case fin := <-finished:
			{
				assert.True(t, fin, "expecting finished to be reached without an error on the error channel")
				return
			}
		}
	}

}

func TestBuilder_Build2(t *testing.T) {
	exists := true
	if _, err := os.Stat("./ocib"); os.IsNotExist(err) {
		exists = false
	}
	assert.False(t, exists, "There should be no context directory (./ocib) after a build has finished executing")
}

func TestBuilder_Pull(t *testing.T) {
}

func TestBuilder_Push(t *testing.T) {
}

func TestBuilder_Login(t *testing.T) {
}

func TestBuilder_Clean(t *testing.T) {
}

func TestBuilder_Purge(t *testing.T) {

}

func (t testClient) ImageBuild(options types.OCIBuildOptions) (types.OCIBuildResponse, error) {
	return types.OCIBuildResponse{}, nil
}

func (t testClient) ImagePull(options types.OCIPullOptions) (types.OCIPullResponse, error) {
	return types.OCIPullResponse{}, nil
}
func (t testClient) ImagePush(options types.OCIPushOptions) (types.OCIPushResponse, error) {
	return types.OCIPushResponse{}, nil
}
func (t testClient) ImageRemove(options types.OCIRemoveOptions) (types.OCIRemoveResponse, error) {
	return types.OCIRemoveResponse{}, nil
}
func (t testClient) RegistryLogin(options types.OCILoginOptions) (types.OCILoginResponse, error) {
	return types.OCILoginResponse{}, nil
}

func (t testClient) GenerateAuthRegistryString(auth dockertypes.AuthConfig) string {
	return ""
}

type testClient struct {
}
