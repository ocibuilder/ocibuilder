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

	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

// S3BuildContextReader implements the context reader for the S3 storage
type S3BuildContextReader struct {
	// buildContext holds the S3 bucket configuration
	buildContext *v1alpha1.S3Context
	// k8sClient is Kubernetes API client
	k8sClient kubernetes.Interface
}

// newSession returns a S3BuildContextReader session
func (contextReader *S3BuildContextReader) newSession(accessKey, secretKey string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Endpoint: &contextReader.buildContext.Endpoint,
		Region:   &contextReader.buildContext.Region,
		Credentials: awscreds.NewStaticCredentialsFromCreds(awscreds.Value{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		}),
		DisableSSL: &contextReader.buildContext.Insecure,
	})
}

// Read reads the context stored on S3BuildContextReader
func (contextReader *S3BuildContextReader) Read() (string, error) {
	accessKey, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.AccessKey)
	if err != nil {
		return "", err
	}
	secretKey, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.SecretKey)
	if err != nil {
		return "", err
	}
	awsSession, err := contextReader.newSession(accessKey, secretKey)
	if err != nil {
		return "", err
	}
	s3Downloader := s3manager.NewDownloader(awsSession)
	contextFilePath := fmt.Sprintf("%s/%s", common.ContextDirectory, common.ContextFile)
	if err := os.MkdirAll(common.ContextDirectory, 0750); err != nil {
		return "", err
	}
	contextFile, err := os.Create(contextFilePath)
	if err != nil {
		return "", err
	}
	if _, err := s3Downloader.Download(contextFile, &awss3.GetObjectInput{
		Bucket: aws.String(contextReader.buildContext.Bucket.Name),
		Key:    aws.String(contextReader.buildContext.Bucket.Key),
	}); err != nil {
		return "", err
	}
	if err := common.UntarFile(contextFilePath, common.ContextDirectoryUncompressed); err != nil {
		return "", err
	}
	return common.ContextDirectoryUncompressed, nil
}

// NewS3BuildContextReader returns a new build context reader for S3
func NewS3BuildContextReader(buildContext *v1alpha1.S3Context, k8sClient kubernetes.Interface) BuildContextReader {
	return &S3BuildContextReader{
		buildContext,
		k8sClient,
	}
}
