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

package ocibuilder

const (
	// Group is the API Group
	Group string = "ocibuilder"

	// Kind is the kind constant for the sensor
	Kind string = "BuildSpecification"
	// Singular is the singular constant for ocibuilder
	Singular string = "buildspecification"
	// Plural is the plural constant for ocibuilder
	Plural string = "buildspecifications"
	// FullName is the full name constant for the sensor
	FullName string = Plural + "." + Group
)
