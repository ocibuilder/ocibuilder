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

package common

import (
	"archive/zip"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientConfig return rest config, if path not specified, assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// UntarFile un-zip/tar a file
func UntarFile(input, destination string) error {
	reader, err := zip.OpenReader(input)
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		destinationFilePath := filepath.Join(destination, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destinationFilePath, file.Mode()); err != nil {
				return err
			}
			continue
		}
		outputFile, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outputFile.Close()

		if _, err = io.Copy(outputFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

// ReadCredentials reads the credentials
func ReadCredentials(client kubernetes.Interface, creds *v1alpha1.Credentials) (string, error) {
	if creds.Plain != "" {
		return creds.Plain, nil
	}
	if creds.Env != "" {
		value, ok := os.LookupEnv(creds.Env)
		if !ok {
			return "", errors.Errorf("environment variable %s for the credentials not found", creds.Env)
		}
		return value, nil
	}
	if creds.KubeSecret != nil {
		if client == nil {
			return "", errors.New("kubernetes client is not initialized")
		}
		secret, err := client.CoreV1().Secrets(creds.KubeSecret.Namespace).Get(creds.KubeSecret.Secret.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		value, ok := secret.Data[creds.KubeSecret.Secret.Key]
		if !ok {
			return "", errors.Errorf("key %s not found in secret %s", creds.KubeSecret.Secret.Key, creds.KubeSecret.Secret.Name)
		}
		return string(value), nil
	}
	return "", errors.New("unknown credentials format")
}
