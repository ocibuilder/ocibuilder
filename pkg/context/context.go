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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/common"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// BuildContextReader enables reading build context from a store
type BuildContextReader interface {
	Read() (string, error)
}

// GetBuildContextReader returns a build context based on the store
func GetBuildContextReader(buildContext *v1alpha1.BuildContext, k8sConfigPath string) (BuildContextReader, error) {
	var k8sClient kubernetes.Interface
	kubeConfig, err := util.GetClientConfig(k8sConfigPath)
	if err == nil {
		k8sClient, err = kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return nil, err
		}
	}
	if buildContext.AliyunOSSContext != nil {
		return NewAliyunOSSBuildContextReader(buildContext.AliyunOSSContext, k8sClient), nil
	}
	if buildContext.AzureBlobContext != nil {
		return NewAzureBlobBuildContextReader(buildContext.AzureBlobContext, k8sClient), nil
	}
	if buildContext.GitContext != nil {
		return NewGitBuildContextReader(buildContext.GitContext, k8sClient), nil
	}
	if buildContext.GCSContext != nil {
		return NewGCSBuildContextReader(buildContext.GCSContext, k8sClient), nil
	}
	if buildContext.LocalContext != nil {
		return NewLocalBuildContextReader(buildContext.LocalContext), nil
	}
	if buildContext.S3Context != nil {
		return NewS3BuildContextReader(buildContext.S3Context, k8sClient), nil
	}
	return nil, errors.New("unknown build context")
}

// InjectDockerfile embeds the generated ocibuilder dockerfile into your build context tar
// looking in /ocib/context/context.tar.gz
func InjectDockerfile(contextPath string, dockerfilePath string) error {

	contextDirectoryPath := fmt.Sprintf("%s%s", contextPath, common.ContextDirectory)
	contextTar := fmt.Sprintf("%s%s", contextDirectoryPath, common.ContextFile)

	if err := util.UntarFile(contextTar, contextDirectoryPath); err != nil {
		return errors.Wrap(err, "error extracting original context file at: "+contextTar)
	}

	if err := os.Remove(contextTar); err != nil {
		return err
	}

	if err := os.Rename(dockerfilePath, fmt.Sprintf("%s%s", contextDirectoryPath, filepath.Base(dockerfilePath))); err != nil {
		return errors.Wrap(err, "error attempting to move Dockerfile to new context directory")
	}

	if err := util.TarFile([]string{contextDirectoryPath}, contextDirectoryPath+common.ContextFile); err != nil {
		return errors.Wrap(err, "error tarring new directory with injected Dockerfile")
	}

	return nil
}

// ExcludeIgnored excludes any explicitly ignored files or directories from the
// build context
func ExcludeIgnored(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	ignored, err := ioutil.ReadFile(directory + "/.dockerignore")
	if err != nil {
		return nil, err
	}

	ignoredPaths := strings.Split(string(ignored), "\n")
	var contextPaths []string
	for _, f := range files {
		if !utils.Exists(f.Name(), ignoredPaths) {
			contextPaths = append(contextPaths, f.Name())
		}
	}

	return contextPaths, nil
}
