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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"github.com/ocibuilder/ocibuilder/testing/dummy"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	builder := Builder{
		Logger:     util.GetLogger(true),
		Client:     testClient{},
		Provenance: []*v1alpha1.BuildProvenance{},
	}

	res := make(chan v1alpha1.OCIBuildResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
		close(finished)
	}()

	go builder.Build(dummy.Spec, res, errChan, finished)

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

func (t testClient) ImageBuild(options v1alpha1.OCIBuildOptions) (v1alpha1.OCIBuildResponse, error) {
	body := ioutil.NopCloser(strings.NewReader("image build response"))
	return v1alpha1.OCIBuildResponse{
		ImageBuildResponse: types.ImageBuildResponse{
			Body:   body,
			OSType: "",
		},
	}, nil
}

func (t testClient) ImagePull(options v1alpha1.OCIPullOptions) (v1alpha1.OCIPullResponse, error) {
	return v1alpha1.OCIPullResponse{}, nil
}
func (t testClient) ImagePush(options v1alpha1.OCIPushOptions) (v1alpha1.OCIPushResponse, error) {
	return v1alpha1.OCIPushResponse{}, nil
}
func (t testClient) ImageRemove(options v1alpha1.OCIRemoveOptions) (v1alpha1.OCIRemoveResponse, error) {
	return v1alpha1.OCIRemoveResponse{}, nil
}
func (t testClient) ImageInspect(imageId string) (types.ImageInspect, error) {
	return types.ImageInspect{}, nil
}
func (t testClient) ImageHistory(imageId string) ([]image.HistoryResponseItem, error) {
	return nil, nil
}
func (t testClient) RegistryLogin(options v1alpha1.OCILoginOptions) (v1alpha1.OCILoginResponse, error) {
	return v1alpha1.OCILoginResponse{}, nil
}
func (t testClient) GenerateAuthRegistryString(auth types.AuthConfig) string {
	return ""
}

type testClient struct {
}
