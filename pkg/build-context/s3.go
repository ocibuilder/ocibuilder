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

	"github.com/aws/aws-sdk-go/aws"
	awscreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// S3BuildContextReader implements the context reader for the S3 storage
type S3BuildContextReader struct {
	// buildContext holds the S3 bucket configuration
	buildContext *v1alpha1.S3Context
	// k8sClient is Kubernetes API client
	k8sClient kubernetes.Interface
}

// credentials holds the S3 credentials
type credentials struct {
	accesskey string
	secretkey string
}

// readSecret reads S3BuildContextReader credentials stored in a Kubernetes secret
func (contextReader *S3BuildContextReader) readFromSecret() (*credentials, error) {
	access, err := common.ReadFromSecret(contextReader.k8sClient, contextReader.buildContext.K8sCreds.Namespace, contextReader.buildContext.K8sCreds.AccessKey)
	if err != nil {
		return nil, err
	}
	secret, err := common.ReadFromSecret(contextReader.k8sClient, contextReader.buildContext.K8sCreds.Namespace, contextReader.buildContext.K8sCreds.SecretKey)
	if err != nil {
		return nil, err
	}
	return &credentials{
		accesskey: string(access),
		secretkey: string(secret),
	}, nil
}

// readFromEnv reads S3BuildContextReader credentials from environment variables
func (contextReader *S3BuildContextReader) readFromEnv() (*credentials, error) {
	access, ok := os.LookupEnv(contextReader.buildContext.EnvVarCreds.EnvVarAccessKey)
	if !ok {
		return nil, errors.New("access key environment variable not found")
	}
	secret, ok := os.LookupEnv(contextReader.buildContext.EnvVarCreds.EnvVarSecretKey)
	if !ok {
		return nil, errors.New("secret key environment variable not found")
	}
	return &credentials{
		accesskey: access,
		secretkey: secret,
	}, nil
}

// getCredentials returns the S3BuildContextReader credentials based on the type of secret store
func (contextReader *S3BuildContextReader) getCredentials() (*credentials, error) {
	if contextReader.buildContext.PlainCreds != nil {
		return &credentials{
			accesskey: contextReader.buildContext.PlainCreds.AccessKey,
			secretkey: contextReader.buildContext.PlainCreds.SecretKey,
		}, nil
	}
	if contextReader.buildContext.EnvVarCreds != nil {
		return contextReader.readFromEnv()
	}
	if contextReader.buildContext.K8sCreds != nil {
		return contextReader.readFromSecret()
	}
	return nil, errors.New("contextReader credentials are not provided")
}

// newSession returns a S3BuildContextReader session
func (contextReader *S3BuildContextReader) newSession(creds *credentials) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Endpoint: &contextReader.buildContext.Endpoint,
		Region:   &contextReader.buildContext.Region,
		Credentials: awscreds.NewStaticCredentialsFromCreds(awscreds.Value{
			AccessKeyID:     creds.accesskey,
			SecretAccessKey: creds.secretkey,
		}),
		DisableSSL: &contextReader.buildContext.Insecure,
	})
}

// Read reads the context stored on S3BuildContextReader
func (contextReader *S3BuildContextReader) Read() error {
	creds, err := contextReader.getCredentials()
	if err != nil {
		return err
	}
	awsSession, err := contextReader.newSession(creds)
	if err != nil {
		return err
	}

	s3Downloader := s3manager.NewDownloader(awsSession)
	contextFilePath := fmt.Sprintf("%s/%s", common.ContextDirectory, common.ContextFile)

	if err := os.MkdirAll(common.ContextDirectory, 0750); err != nil {
		return err
	}
	contextFile, err := os.Create(contextFilePath)
	if err != nil {
		return err
	}
	if _, err := s3Downloader.Download(contextFile, &awss3.GetObjectInput{
		Bucket: aws.String(contextReader.buildContext.Bucket.Name),
		Key:    aws.String(contextReader.buildContext.Bucket.Key),
	}); err != nil {
		return err
	}
	if err := common.UntarFile(contextFilePath, common.ContextDirectoryUncompressed); err != nil {
		return err
	}
	return nil
}
