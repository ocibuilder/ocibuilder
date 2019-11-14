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

package build_context

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// AliyunOSSBuildContextReader implements BuildContextReader to fetch build context stored on Aliyun OSS
type AliyunOSSBuildContextReader struct {
	// buildContext stored on OSS
	buildContext *v1alpha1.AliyunOSSContext
	// k8sClient is Kubernetes client
	k8sClient kubernetes.Interface
}

// Read reads and stores build context from OSS
func (contextReader *AliyunOSSBuildContextReader) Read() (string, error) {
	accessId, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.AccessId)
	if err != nil {
		return "", err
	}
	accessSecret, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.AccessSecret)
	if err != nil {
		return "", err
	}
	client, err := oss.New(contextReader.buildContext.Endpoint, accessId, accessSecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(contextReader.buildContext.Bucket.Name)
	if err != nil {
		return "", err
	}
	contextFilePath := fmt.Sprintf("%s/%s", common.ContextDirectory, common.ContextFile)
	if err := os.MkdirAll(common.ContextDirectory, 0750); err != nil {
		return "", err
	}
	if _, err := os.Create(contextFilePath); err != nil {
		return "", err
	}
	if err := bucket.GetObjectToFile(contextReader.buildContext.Bucket.Key, contextFilePath); err != nil {
		return "", err
	}
	if err := common.UntarFile(contextFilePath, common.ContextDirectoryUncompressed); err != nil {
		return "", err
	}
	return common.ContextDirectoryUncompressed, nil
}
