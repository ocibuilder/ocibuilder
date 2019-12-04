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

package context

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestLocalBuildContextReader_Read(t *testing.T) {

	reader, err := GetBuildContextReader(buildContext, "")
	assert.Equal(t, nil, err)

	path, err := reader.Read()
	assert.Equal(t, nil, err)
	assert.Equal(t, TEST_SERVICE_PATH, path)

	err = common.UntarFile(TEST_SERVICE_PATH+"/ocib/context/context.tar.gz", TEST_SERVICE_PATH+"/unpacked")
	assert.Equal(t, nil, err)

	files, err := ioutil.ReadDir(TEST_SERVICE_PATH + "/unpacked/")
	assert.Equal(t, nil, err)

	var actualFileNames []string
	for _, file := range files {
		actualFileNames = append(actualFileNames, file.Name())
	}
	assert.Equal(t, expectedFileNames, actualFileNames)

	err = os.RemoveAll(TEST_SERVICE_PATH + "/ocib")
	assert.Equal(t, nil, err)

	err = os.RemoveAll(TEST_SERVICE_PATH + "/unpacked")
	assert.Equal(t, nil, err)
}

const TEST_SERVICE_PATH = "../../testing/e2e/resources/go-test-service"

var buildContext = &v1alpha1.BuildContext{
	LocalContext: &v1alpha1.LocalContext{
		ContextPath: TEST_SERVICE_PATH,
	},
}

var expectedFileNames = []string{"main.go", "ocibuilder.yaml"}
