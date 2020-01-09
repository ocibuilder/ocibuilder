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
	Metadata *v1alpha1.BuildMetadata
	Logger   *logrus.Logger
	Store    store.MetadataStore
	// records holds all the records that have been parsed ready to push
	records []*store.Record
}

func (m MetadataWriter) Write() error {
	if err := m.Store.Write(m.records...); err != nil {
		return err
	}
	return nil
}

func (m *MetadataWriter) ParseMetadata(buildResponse io.ReadCloser) error {

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(buildResponse); err != nil {
		return err
	}

	responseOutput := strings.Split(buf.String(), "\n")
	var layerIds []string
	var layerInfo []types.ImageLayer

	for i, line := range responseOutput {

		if (strings.Contains(line, "Step") && i != 0) || strings.Contains(line, "Successfully built") {
			// Retrieve the layer digest from the previous line
			sep := strings.Split(responseOutput[i-1], " ")
			layerIds = append(layerIds, sep[len(sep)-1])
		}

		if strings.Contains(line, "Step") {
			// Retrieve the step line including only the command and the args to the command
			cmdLine := strings.Split(responseOutput[i], " : ")[1]
			// Separate the specific command being executed for each layer
			cmd := types.LayerDirective(strings.Split(cmdLine, " ")[0])
			// Join the remaining args into a single string to be stored
			args := strings.Join(strings.Split(cmdLine, " ")[1:], " ")

			layerInfo = append(layerInfo, types.ImageLayer{
				Directive: &cmd,
				Arguments: args,
			})
		}

	}

	if len(layerIds) == 0 {
		return errors.New("no image ids found in image response")
	}

	imageId := layerIds[len(layerIds)-1]
	m.createAttestation(imageId)

	imageDetailMetdata := types.V1beta1imageDetails{
		DerivedImage: &types.ImageDerived{
			Fingerprint: &types.ImageFingerprint{
				V1Name: imageId,
				V2Blob: layerIds[:len(layerIds)-1],
			},
			LayerInfo: layerInfo,
		},
	}
	record := store.Record{DerivedImage: &imageDetailMetdata}
	m.records = append(m.records, &record)

	return nil

}

func (m *MetadataWriter) createAttestation(digest string) {

}

func NewMetadataWriter(logger *logrus.Logger, metadataSpec *v1alpha1.BuildMetadata) MetadataWriter {
	var metadataStore store.MetadataStore

	if metadataSpec.StoreConfig.Grafeas != nil {

		config := types.Configuration{
			BasePath:   metadataSpec.Hostname,
			HTTPClient: &http.Client{},
		}

		metadataStore = grafeas.NewStore(&config, metadataSpec.StoreConfig.Grafeas, logger)
	}

	return MetadataWriter{
		Logger:   logger,
		Store:    metadataStore,
		Metadata: metadataSpec,
	}
}
