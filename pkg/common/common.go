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

package common

import "github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder"

// Controller environment variables
const (
	// EnvVarControllerConfigMap is the name of the configmap to use for the controller
	EnvVarControllerConfigMap = "CONTROLLER_CONFIG_MAP"
	EnvVarKubeConfig          = "KUBE_CONFIG"
)

// Controller labels
const (
	//LabelKeyControllerInstanceID is the label which allows to separate application among multiple running ocibuilder controllers.
	LabelKeyControllerInstanceID = ocibuilder.FullName + "/controller-instanceid"
	// LabelKeyPhase is a label applied to indicate the current phase of the builder (for filtering purposes)
	LabelKeyPhase = ocibuilder.FullName + "/phase"
	// LabelKeyComplete is the label to mark builders as complete
	LabelKeyComplete = ocibuilder.FullName + "/complete"
	// LabelOCIBuilderName is the label to indicate the name of an ocibuilder object
	LabelOCIBuilderName = "ocibuilder-name"
)

// Miscellaneous constants for controller
const (
	// ControllerConfigMapKey is the key in the configmap to retrieve ocibuilder controller configuration from.
	// Content encoding is expected to be YAML.
	ControllerConfigMapKey = "config"
)

// OCIBuilder resource labels
const (
	// LabelName refers to label of the ocibuilder resource name
	LabelName = "name"
	// LabelNamespace is the label to indicate K8s namespace
	LabelNamespace = "namespace"
	// LabelPhase is the label to indicate the phase of the ocibuilder resourc
	LabelPhase = "phase"
)

// logger labels
const (
	LabelVersion = "version"
)

// Default image registry
const (
	DefaultImageRegistry = "docker.io"
)

// Build context constants
const (
	// ContextDirectory holds the ocibuilder context
	ContextDirectory = "/ocib/context/"
	// ContextFile contains the compressed build context
	ContextFile = "context.tar.gz"
	// ContextDirectoryUncompressed contains the uncompressed build context
	ContextDirectoryUncompressed = "/ocibuilder/context/uncompressed/"
)
