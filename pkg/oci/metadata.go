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
	"errors"
	"net/http"

	types "github.com/artbegolli/grafeas"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
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

func (m *MetadataWriter) ParseMetadata(imageName string, cli v1alpha1.BuilderClient) error {

	if _, ok := cli.(docker.Client); !ok {
		return errors.New("writing metadata not currently supported for use with buildah")
	}

	inspectResponse, err := cli.ImageInspect(imageName)
	if err != nil {
		return err
	}

	historyResponse, err := cli.ImageHistory(imageName)
	if err != nil {
		return err
	}

	var layerIds []string
	var layerInfo []types.ImageLayer
	for _, r := range historyResponse {
		layerIds = append(layerIds, r.ID)
		layerInfo = append(layerInfo, types.ImageLayer{
			Arguments: r.CreatedBy,
		})
	}

	record := store.Record{
		DerivedImage: &types.V1beta1imageDetails{
			DerivedImage: &types.ImageDerived{
				Fingerprint: &types.ImageFingerprint{
					V1Name: inspectResponse.ID,
					V2Blob: layerIds,
				},
				LayerInfo: layerInfo,
			},
		},
	}
	m.records = append(m.records, &record)
	m.createAttestation(inspectResponse.RepoDigests[0])

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