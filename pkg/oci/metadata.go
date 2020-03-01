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
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/ocibuilder/gofeas"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/crypto"
	"github.com/ocibuilder/ocibuilder/pkg/store"
	"github.com/ocibuilder/ocibuilder/pkg/store/grafeas"
	"github.com/sirupsen/logrus"
)

type MetadataWriter struct {
	Metadata *v1alpha1.Metadata
	Logger   *logrus.Logger
	Store    store.MetadataStore
	// records holds all the records that have been parsed ready to push
	records []*store.Record
	// resource holds the resource ID
	resource string
}

func (m MetadataWriter) Write() error {
	if err := m.Store.Write(m.records...); err != nil {
		return err
	}
	return nil
}

func (m *MetadataWriter) ParseMetadata(imageName string, cli v1alpha1.BuilderClient, provenance *v1alpha1.BuildProvenance) error {
	log := m.Logger

	if _, ok := cli.(buildah.Client); ok {
		return errors.New("writing metadata is currently only supported for use with docker")
	}

	log.Debugln("conducting image inspect")
	inspectResponse, err := cli.ImageInspect(imageName)
	if err != nil || inspectResponse.ID == "" {
		log.Errorln("error in inspecting image or no response ID returned - cannot push metadata")
		return err
	}
	log.WithField("responseId", inspectResponse.ID).Debugln("inspect response ID")

	log.Debugln("conducting image history")
	historyResponse, err := cli.ImageHistory(imageName)
	if err != nil || historyResponse[0].ID == "" {
		log.Errorln("error in request image history or no history ID response returned - cannot push metadata")
		return err
	}

	var layerIds []string
	var layerInfo []gofeas.ImageLayer
	for _, r := range historyResponse {
		layerIds = append(layerIds, r.ID)
		layerInfo = append(layerInfo, parseCreatedBy(r.CreatedBy))
	}

	// Currently we are signing ImageID instead of digests - digests are generated when a local manifest is created pre-push
	digest := inspectResponse.ID
	m.resource = fmt.Sprintf("%s@%s", provenance.Name, digest)

	layerIds = append(layerIds, inspectResponse.ID)

	if len(m.Metadata.Data) == 0 {
		log.Errorln("no requested metadata types defined in the specification - ensure that metadata types are defined in the data field")
		return errors.New("no requested metadata types found")
	}

	for _, requestedMetadata := range m.Metadata.Data {

		switch requestedMetadata {

		case v1alpha1.Build:
			{
				log.Debugln("found request for build metadata, creating build record")
				buildRecord := m.createBuildRecord(digest, provenance)
				m.records = append(m.records, &buildRecord)
			}
		case v1alpha1.Attestation:
			{
				log.Debugln("found request for attestation metadata, creating attestation record")
				attestationRecord, err := m.createAttestationRecord(digest)
				if err != nil {
					log.Errorln("error creating attestation record")
					return err
				}
				m.records = append(m.records, &attestationRecord)
			}
		case v1alpha1.Image:
			{
				log.Debugln("found request for image metadata, creating image record")
				derivedImageRecord := m.createDerivedImageRecord(inspectResponse.ID, layerIds, layerInfo)
				m.records = append(m.records, &derivedImageRecord)
			}
		default:
			{
				log.Warnf("the metadata type %s defined is not valid\n", requestedMetadata)
			}
		}

	}

	return nil

}

func (m *MetadataWriter) createBuildRecord(digest string, buildProvenance *v1alpha1.BuildProvenance) store.Record {
	derivedBuildRecord := store.Record{
		Resource: m.resource,
		Build: &gofeas.V1beta1buildDetails{
			Provenance: &gofeas.ProvenanceBuildProvenance{
				Id:        uuid.New().String(),
				ProjectId: m.Metadata.StoreConfig.Grafeas.Project,
				BuiltArtifacts: []gofeas.ProvenanceArtifact{{
					Checksum: digest,
					Id:       m.resource,
					Names:    []string{fmt.Sprintf("%s:%s", buildProvenance.Name, buildProvenance.Tag)},
				}},
				StartTime:  buildProvenance.StartTime,
				EndTime:    buildProvenance.EndTime,
				CreateTime: buildProvenance.EndTime,
				Creator:    buildProvenance.Creator,
				SourceProvenance: &gofeas.ProvenanceSource{
					ArtifactStorageSourceUri: buildProvenance.Source,
				},
			},
		},
	}
	return derivedBuildRecord
}

func (m *MetadataWriter) createDerivedImageRecord(imageId string, layerIds []string, layerInfo []gofeas.ImageLayer) store.Record {
	derivedImageRecord := store.Record{
		Resource: m.resource,
		DerivedImage: &gofeas.V1beta1imageDetails{
			DerivedImage: &gofeas.ImageDerived{
				Fingerprint: &gofeas.ImageFingerprint{
					V1Name: imageId,
					V2Blob: layerIds,
				},
				LayerInfo: layerInfo,
			},
		},
	}
	return derivedImageRecord
}

func (m *MetadataWriter) createAttestationRecord(digest string) (store.Record, error) {

	if m.Metadata.Key == nil {
		return store.Record{}, errors.New("no signing key has been defined for image attestation")
	}

	privKey, pubKey, err := crypto.ValidateKeysPacket(m.Metadata.Key)
	if err != nil {
		return store.Record{}, err
	}
	e := crypto.CreateEntityFromKeys(privKey, pubKey)
	id, sig, err := crypto.SignDigest(digest, m.Metadata.Key.Passphrase, e)
	if err != nil {
		return store.Record{}, err
	}

	record := store.Record{
		Resource: m.resource,
		Attestation: &gofeas.V1beta1attestationDetails{
			Attestation: &gofeas.AttestationAttestation{
				PgpSignedAttestation: &gofeas.AttestationPgpSignedAttestation{
					Signature: sig,
					PgpKeyId:  id,
				},
			},
		},
	}

	return record, nil
}

func NewMetadataWriter(logger *logrus.Logger, metadataSpec *v1alpha1.Metadata) MetadataWriter {
	var metadataStore store.MetadataStore

	if metadataSpec.StoreConfig.Grafeas != nil {

		config := gofeas.Configuration{
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

func parseCreatedBy(createdBy string) gofeas.ImageLayer {
	createdBy = strings.TrimPrefix(createdBy, "/bin/sh -c ")

	if strings.HasPrefix(createdBy, "#(nop)") {
		createdBySplit := strings.Fields(createdBy)

		directive := gofeas.LayerDirective(createdBySplit[1])
		return gofeas.ImageLayer{
			Directive: &directive,
			Arguments: strings.Join(createdBySplit[2:], " "),
		}
	}

	directive := gofeas.RUN_LayerDirective
	return gofeas.ImageLayer{
		Directive: &directive,
		Arguments: createdBy,
	}
}
