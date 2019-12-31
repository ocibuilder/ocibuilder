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
	ctx "context"
	"time"

	"github.com/artbegolli/grafeas"
	"github.com/ocibuilder/ocibuilder/pkg/store"
)

type graf struct {
	// Client is the grafeas API client
	Client *grafeas.APIClient
	// Project is the metadata project for grafeas image build metadata
	Project string
}

func (g *graf) List() ([]*store.Record, error) {
	return nil, nil
}

func (g *graf) Read(key ...string) ([]*store.Record, error) {
	return nil, nil
}

func (g *graf) Write(rec ...*store.Record) error {

	requestMap := make(map[string]grafeas.V1beta1Occurrence)
	for _, r := range rec {
		requestMap[r.Key] = grafeas.V1beta1Occurrence{
			Name:          "",
			Resource:      nil,
			NoteName:      "",
			Kind:          nil,
			Remediation:   "",
			CreateTime:    time.Time{},
			UpdateTime:    time.Time{},
			Vulnerability: nil,
			Build:         nil,
			DerivedImage:  nil,
			Installation:  nil,
			Deployment:    nil,
			Discovered:    nil,
			Attestation:   nil,
		}
	}

	batchNotesRequest := grafeas.V1beta1BatchCreateNotesRequest{}

	_, _, err := g.Client.GrafeasV1Beta1Api.BatchCreateNotes(ctx.Background(), g.Project, batchNotesRequest)
	if err != nil {
		return err
	}

	return nil
}

func (g *graf) Delete(key ...string) error {
	return nil
}

func NewStore(project string, configuration *grafeas.Configuration) store.Store {
	cli := grafeas.NewAPIClient(configuration)

	return &graf{
		Project: project,
		Client:  cli,
	}
}
