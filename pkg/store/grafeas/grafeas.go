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
	"fmt"
	"net/http"
	"time"

	"github.com/artbegolli/grafeas"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/store"
	"github.com/sirupsen/logrus"
)

type graf struct {
	// Client is the grafeas API client
	Client *grafeas.APIClient
	// Options stores the options for pushing to Grafeas
	Options v1alpha1.Grafeas
	// Logger is the logger
	Logger *logrus.Logger
}

// Write writes to the Grafeas metadata store. It is a variadic function
// and takes in a number of Records.
// The records are parsed as follows
func (g *graf) Write(rec ...*store.Record) error {

	var occurrenceRequests []grafeas.V1beta1Occurrence
	for _, r := range rec {

		occ := grafeas.V1beta1Occurrence{
			Name:       "",
			Resource:   nil,
			NoteName:   "",
			Kind:       nil,
			CreateTime: time.Now(),
		}

		if r.Build != nil {
			occ.Build = r.Build
			continue
		}

		if r.DerivedImage != nil {
			occ.DerivedImage = r.DerivedImage
			continue
		}

		if r.Attestation != nil {
			occ.Attestation = r.Attestation
			continue
		}

	}

	req := grafeas.V1beta1BatchCreateOccurrencesRequest{
		Parent:      "",
		Occurrences: occurrenceRequests,
	}

	res, httpRes, err := g.Client.GrafeasV1Beta1Api.BatchCreateOccurrences(ctx.Background(), g.Options.Project, req)

	if err != nil {
		return err
	}

	if httpRes.StatusCode != http.StatusOK {
		return fmt.Errorf("error making write request to grafeas - returned with status code %s", httpRes.Status)
	}

	for _, occurrenceResponse := range res.Occurrences {
		g.Logger.WithFields(logrus.Fields{
			"name":        occurrenceResponse.Name,
			"create_time": occurrenceResponse.CreateTime,
			"kind":        occurrenceResponse.Kind,
		}).Debugln("finished pushing metadata to Grafeas")
	}

	return nil
}

func NewStore(configuration *grafeas.Configuration, options v1alpha1.Grafeas, logger *logrus.Logger) store.MetaStore {
	cli := grafeas.NewAPIClient(configuration)

	return &graf{
		Client:  cli,
		Options: options,
		Logger:  logger,
	}
}
