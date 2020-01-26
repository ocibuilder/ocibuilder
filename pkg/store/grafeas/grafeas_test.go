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

package grafeas

import (
	"net/http"
	"testing"

	"github.com/ocibuilder/gofeas"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/store"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestGraf_Write(t *testing.T) {
	var testStore = graf{
		Client:  testClient{T: t},
		Options: options,
		Logger:  util.Logger,
	}

	err := testStore.Write(record)
	assert.Equal(t, nil, err)
}

var record = &store.Record{
	Resource: "random-occ-resource",
	Attestation: &gofeas.V1beta1attestationDetails{
		Attestation: &gofeas.AttestationAttestation{
			PgpSignedAttestation: &gofeas.AttestationPgpSignedAttestation{
				Signature: "this-is-a-signature",
				PgpKeyId:  "1",
			},
		},
	},
}

var options = &v1alpha1.Grafeas{
	Project:  "image-signing",
	NoteName: "projects/image-signing/notes/production",
	Resource: "random-resource",
}

func (t testClient) BatchCreateOccurrences(ctx context.Context, parent string, body gofeas.V1beta1BatchCreateOccurrencesRequest) (gofeas.V1beta1BatchCreateOccurrencesResponse, *http.Response, error) {
	assert.Equal(t.T, "projects/image-signing", parent)
	assert.Equal(t.T, "random-occ-resource", body.Occurrences[0].Resource.Uri)
	assert.Equal(t.T, body.Occurrences[0].Attestation, record.Attestation)
	return gofeas.V1beta1BatchCreateOccurrencesResponse{}, nil, nil
}

type testClient struct {
	gofeas.APIClient
	T *testing.T
}

var expectedRequest = gofeas.V1beta1BatchCreateOccurrencesRequest{
	Parent: "projects/image-signing",
	Occurrences: []gofeas.V1beta1Occurrence{{
		Resource: &gofeas.V1beta1Resource{
			Uri: "random-resource",
		},
		NoteName:    "projects/image-signing/notes/production",
		Attestation: record.Attestation,
	}},
}
