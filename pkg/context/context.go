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
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
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
	kubeConfig, err := common.GetClientConfig(k8sConfigPath)
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

// InjectDockerfile embeds the generated ocibuilder dockerfile into your build context
func InjectDockerfile(contextPath string, dockerfilePath string) error {
	destination := fmt.Sprintf("%s%s", strings.Replace(contextPath, common.ContextFile, "", -1), "context/")
	if err := common.UntarFile(contextPath, destination); err != nil {
		logrus.WithError(err).Infoln("error untarring")
		return err
	}

	if err := os.Remove(contextPath); err != nil {
		logrus.WithError(err).Infoln("error removing tar")
		return err
	}

	if err := os.Rename(dockerfilePath, fmt.Sprintf("%s%s", destination, filepath.Base(dockerfilePath))); err != nil {
		logrus.WithError(err).Infoln("error moving dockerfile")
		return err
	}

	if err := common.TarFile(destination, contextPath); err != nil {
		logrus.WithError(err).Infoln("error tarring")
		return err
	}

	return nil
}
