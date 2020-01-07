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
	"io"
	"net/http"

	client "github.com/artbegolli/grafeas"
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

func (m MetadataWriter) Write() {

}

func (m *MetadataWriter) ParseResponseMetadata(buildResponse io.ReadCloser) {

}

func NewMetadataWriter(logger *logrus.Logger, metadataSpec v1alpha1.BuildMetadata) MetadataWriter {
	var metaStore store.MetaStore

	if metadataSpec.Store.Grafeas != nil {

		config := client.Configuration{
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
