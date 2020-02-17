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

package utils

import (
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = common.GetLogger(false)

// OutputJson streams and formats the output to stdout from returned ReadClosers by docker
// commands.
func OutputJson(output io.ReadCloser) error {

	termFd, isTerm := term.GetFdInfo(os.Stdout)

	err := jsonmessage.DisplayJSONMessagesStream(
		output,
		os.Stdout,
		termFd,
		isTerm,
		nil,
	)

	if err != nil {
		log.WithError(err).Errorln("failed to get JSON stream")
		return err
	}

	return nil

}

// Output outputs a readcloser to stdout in a stream without formatting.
func Output(stdout io.ReadCloser, stderr io.ReadCloser) error {
	//TODO: error with premature read |0: file already closed when finished reading out, investigate further
	if _, err := io.Copy(os.Stdout, stderr); err != nil {
		log.WithError(err).Debugln("error copying output from stderr to stdout, could impact response output")
	}

	if _, err := io.Copy(os.Stdout, stdout); err != nil {
		log.WithError(err).Debugln("error copying output from stdout to stdout, could impact response output")
	}
	return nil
}

// GetClient returns a OCIBuilder client
func GetClient(builderType string, logger *logrus.Logger) (v1alpha1.BuilderClient, error) {
	switch v1alpha1.Framework(builderType) {
	case v1alpha1.DockerFramework:
		apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create a docker client")
		}
		return docker.Client{
			APIClient: apiClient,
			Logger:    logger,
		}, nil

	case v1alpha1.BuildahFramework:
		return buildah.Client{
			Logger: logger,
		}, nil
	default:
		return nil, errors.Errorf("invalid builder %s, try --builder=docker or --builder=buildah", builderType)
	}
}

// HasDaemon determines if docker daemon is required for given OCIBuilder type
func HasDaemon(builderType string) bool {
	return v1alpha1.Framework(builderType) == v1alpha1.DockerFramework
}
