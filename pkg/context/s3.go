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
	"io"

	corev1 "k8s.io/api/core/v1"
)

// S3EnvCreds holds the references to environment variables that holds access and secret keys
type S3EnvCreds struct {
	EnvVarAccessKey string `json:"envVarKey" protobuf:"bytes,1,name=envVarAccessKey"`
	EnvVarSecretKey string `json:"envVarSecretKey" protobuf:"bytes,2,name=envVarSecretKey"`
}

// S3PlainCreds holds reference to plain text access and secret keys
type S3PlainCreds struct {
	AccessKey string `json:"accessKey" protobuf:"bytes,1,name=accessKey"`
	SecretKey string `json:"secretKey" protobuf:"bytes,2,name=secretKey"`
}

// S3K8sCreds holds reference to K8s secret that holds access and secret keys
type S3K8sCreds struct {
	AccessKey *corev1.SecretKeySelector `json:"accessKey" protobuf:"bytes,1,name=accessKey"`
	SecretKey *corev1.SecretKeySelector `json:"secretKey" protobuf:"bytes,2,name=secretKey"`
}

// S3Bucket contains information to describe an S3 Bucket
type S3Bucket struct {
	Key  string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

// S3Context refers to context stored on S3 bucket to build an image
type S3Context struct {
	Endpoint    string        `json:"endpoint" protobuf:"bytes,1,name=endpoint"`
	Bucket      *S3Bucket     `json:"bucket" protobuf:"bytes,2,name=bucket"`
	Region      string        `json:"region,omitempty" protobuf:"bytes,3,opt,name=region"`
	Insecure    bool          `json:"insecure,omitempty" protobuf:"varint,4,opt,name=insecure"`
	EnvVarCreds *S3EnvCreds   `json:"envVarCreds,omitempty" protobuf:"bytes,5,opt,name=envVarCreds"`
	PlainCreds  *S3PlainCreds `json:"plainCreds,omitempty" protobuf:"bytes,6,opt,name=plainCreds"`
	K8sCreds    *S3K8sCreds   `json:"k8sCreds,omitempty" protobuf:"bytes,7,opt,name=k8sCreds"`
}

func fetchFromSecret()

// Read reads the context stored on S3
func (ctx S3Context) Read() (io.ReadCloser, error) {
	return nil, nil
}
