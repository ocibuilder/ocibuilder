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
	"fmt"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

// LocalBuildContextReader implements BuildContextReader for local build contexts
type LocalBuildContextReader struct {
	buildContext *v1alpha1.LocalContext
}

// NewLocalBuildContextReader returns a local build context reader
func NewLocalBuildContextReader(buildContext *v1alpha1.LocalContext) *LocalBuildContextReader {
	return &LocalBuildContextReader{
		buildContext,
	}
}

// Read reads the build context from the local
func (contextReader *LocalBuildContextReader) Read() (string, error) {
	contextFilePath := fmt.Sprintf("%s%s/%s", contextReader.buildContext.ContextPath, common.ContextDirectory, common.ContextFile)
	if err := common.TarFile(contextReader.buildContext.ContextPath+"/.", contextFilePath); err != nil {
		return "", err
	}
	return contextReader.buildContext.ContextPath, nil
}
