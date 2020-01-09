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

package store

import (
	"time"

	"github.com/artbegolli/grafeas"
)

// MetadataStore is a data storage interface
type MetadataStore interface {
	// Write records
	Write(rec ...*Record) error
}

// Record represents a data record
type Record struct {
	Key    string
	Value  []byte
	Expiry time.Duration
	// Occurrence is the name of a grafeas occurrence to push to grafeas
	Occurrence string
	// Resource is the resource (e.g. fully qualified image name) of the resource metadata is for
	Resource     string
	Build        *grafeas.V1beta1buildDetails
	DerivedImage *grafeas.V1beta1imageDetails
	Attestation  *grafeas.V1beta1attestationDetails
}
