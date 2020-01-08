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
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/common"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"k8s.io/client-go/kubernetes"
)

// AzureBlobBuildContextReader implements BuildContextReader for the context stored on Azure Storage Blob
type AzureBlobBuildContextReader struct {
	// buildContext contains configuration required to read the build context stored on the Azure Storage Blob
	buildContext *v1alpha1.AzureBlobContext
	// k8sClient is the Kubernetes client
	k8sClient kubernetes.Interface
}

// Read reads the build context from Azure Storage Blob and stores it at a preconfigured path
func (contextReader *AzureBlobBuildContextReader) Read() (string, error) {
	accountName, err := util.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.Account)
	if err != nil {
		return "", err
	}
	accessKey, err := util.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.AccessKey)
	if err != nil {
		return "", nil
	}
	urlStr, err := util.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.URL)
	if err != nil {
		return "", err
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accessKey)
	if err != nil {
		return "", err
	}
	blobURL := azblob.NewBlockBlobURL(*parsedURL, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	downloadResponse, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return "", err
	}
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	downloadedData := bytes.Buffer{}
	if _, err = downloadedData.ReadFrom(bodyStream); err != nil {
		return "", err
	}
	var contextBody []byte
	if _, err := bodyStream.Read(contextBody); err != nil {
		return "", err
	}
	contextFilePath := fmt.Sprintf("%s/%s", common.ContextDirectory, common.ContextFile)
	if err := os.MkdirAll(common.ContextDirectory, 0750); err != nil {
		return "", err
	}
	contextFile, err := os.Create(contextFilePath)
	if err != nil {
		return "", err
	}
	if _, err := contextFile.Write(contextBody); err != nil {
		return "", nil
	}
	if err := util.UntarFile(contextFilePath, common.ContextDirectoryUncompressed); err != nil {
		return "", err
	}
	return common.ContextDirectoryUncompressed, nil
}

// NewAzureBlobBuildContextReader returns a new build context reader for Azure Storage Blob
func NewAzureBlobBuildContextReader(buildContext *v1alpha1.AzureBlobContext, k8sClient kubernetes.Interface) *AzureBlobBuildContextReader {
	return &AzureBlobBuildContextReader{
		buildContext,
		k8sClient,
	}
}
