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
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"k8s.io/client-go/kubernetes"
)

// GCSBuildContextReader implements the BuildContextReader to read build context from Google Cloud Storage
type GCSBuildContextReader struct {
	// buildContext contains configuration required to read build context from GCS
	buildContext *v1alpha1.GCSContext
	// k8sClient is a Kubernetes client
	k8sClient kubernetes.Interface
}

// newClient returns the new GCS client based on authentication methods
func (contextReader *GCSBuildContextReader) newClient() (*storage.Client, error) {
	ctx := context.Background()
	if !contextReader.buildContext.AuthRequired {
		return storage.NewClient(ctx, option.WithoutAuthentication(), option.WithEndpoint(contextReader.buildContext.Endpoint))
	}
	if contextReader.buildContext.CredentialsFilePath != "" {
		return storage.NewClient(ctx, option.WithCredentialsFile(contextReader.buildContext.CredentialsFilePath), option.WithEndpoint(contextReader.buildContext.Endpoint))
	}
	if contextReader.buildContext.APIKey != nil {
		apiKey, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.APIKey)
		if err != nil {
			return nil, err
		}
		return storage.NewClient(ctx, option.WithAPIKey(apiKey), option.WithEndpoint(contextReader.buildContext.Endpoint))
	}
	return nil, errors.New("no authentication method provided. If no authentication is required, set the `authRequired` to true")
}

// Read reads the build context from GCS
func (contextReader *GCSBuildContextReader) Read() (string, error) {
	client, err := contextReader.newClient()
	if err != nil {
		return "", err
	}
	reader, err := client.Bucket(contextReader.buildContext.Bucket.Name).Object(contextReader.buildContext.Bucket.Key).NewReader(context.Background())
	if err != nil {
		return "", err
	}
	var contextBody []byte
	if _, err := reader.Read(contextBody); err != nil {
		return "", nil
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
	if err := common.UntarFile(contextFilePath, common.ContextDirectoryUncompressed); err != nil {
		return "", err
	}
	return common.ContextDirectoryUncompressed, nil
}

// NewGCSBuildContextReader returns a new build context reader for GCS
func NewGCSBuildContextReader(buildContext *v1alpha1.GCSContext, k8sClient kubernetes.Interface) BuildContextReader {
	return &GCSBuildContextReader{
		buildContext,
		k8sClient,
	}
}
