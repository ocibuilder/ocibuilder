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

package oci

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	types "github.com/artbegolli/grafeas"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/store"
	"github.com/ocibuilder/ocibuilder/pkg/store/grafeas"
	"github.com/sirupsen/logrus"
)

type MetadataWriter struct {
	Metadata v1alpha1.BuildMetadata
	Logger   *logrus.Logger
	Store    store.MetaStore
}

func (m MetadataWriter) Write() error {
	return nil
}

func (m *MetadataWriter) ParseResponseMetadata(buildResponse io.ReadCloser) error {

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(buildResponse); err != nil {
		return err
	}

	responseOutput := strings.Split(buf.String(), "\n")
	var layerDigests []string
	var layerInfo []types.ImageLayer

	for i, line := range responseOutput {

		if (strings.Contains(line, "Step") && i != 0) || strings.Contains(line, "Successfully built") {
			sep := strings.Split(responseOutput[i-1], " ")
			layerDigests = append(layerDigests, sep[len(sep)-1])
		}

		if strings.Contains(line, "Step") {
			cmdLine := strings.Split(responseOutput[i], " : ")[1]
			cmd := types.LayerDirective(strings.Split(cmdLine, " ")[0])
			args := strings.Join(strings.Split(cmdLine, " ")[1:], " ")

			layerInfo = append(layerInfo, types.ImageLayer{
				Directive: &cmd,
				Arguments: args,
			})
		}

	}

	if len(layerDigests) == 0 {
		return errors.New("no image digests found in image response")
	}

	imageMetadata := types.V1beta1imageDetails{
		DerivedImage: &types.ImageDerived{
			Fingerprint: &types.ImageFingerprint{
				V1Name: layerDigests[len(layerDigests)-1],
				V2Blob: layerDigests[:len(layerDigests)-1],
			},
			LayerInfo: layerInfo,
		},
	}
	fmt.Println(imageMetadata)

	return nil

}

func NewMetadataWriter(logger *logrus.Logger, metadataSpec *v1alpha1.BuildMetadata) MetadataWriter {
	var metaStore store.MetaStore

	if metadataSpec.Store.Grafeas != nil {

		config := types.Configuration{
			BasePath:   metadataSpec.Hostname,
			HTTPClient: &http.Client{},
		}

		metaStore = grafeas.NewStore(&config, metadataSpec.Store.Grafeas, logger)
	}

	return MetadataWriter{
		Logger: logger,
		Store:  metaStore,
	}
}
