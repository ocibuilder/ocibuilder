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

import "github.com/beval/beval/pkg/apis/beval"

// Controller environment variables
const (
	// EnvVarControllerConfigMap is the name of the configmap to use for the controller
	EnvVarControllerConfigMap = "CONTROLLER_CONFIG_MAP"
	EnvVarKubeConfig          = "KUBE_CONFIG"
)

// Controller labels
const (
	//LabelKeyControllerInstanceID is the label which allows to separate application among multiple running beval controllers.
	LabelKeyControllerInstanceID = beval.FullName + "/controller-instanceid"
	// LabelKeyPhase is a label applied to indicate the current phase of the builder (for filtering purposes)
	LabelKeyPhase = beval.FullName + "/phase"
	// LabelKeyComplete is the label to mark builders as complete
	LabelKeyComplete = beval.FullName + "/complete"
	// LabelbevalName is the label to indicate the name of an beval object
	LabelbevalName = "beval-name"
)

// Miscellaneous constants for controller
const (
	// ControllerConfigMapKey is the key in the configmap to retrieve beval controller configuration from.
	// Content encoding is expected to be YAML.
	ControllerConfigMapKey = "config"
)

// beval resource labels
const (
	// LabelName refers to label of the beval resource name
	LabelName = "name"
	// LabelNamespace is the label to indicate K8s namespace
	LabelNamespace = "namespace"
	// LabelPhase is the label to indicate the phase of the beval resourc
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
	// ContextDirectory holds the beval context
	ContextDirectory = "/ocib/context/"
	// ContextFile contains the compressed build context
	ContextFile = "context.tar.gz"
	// ContextDirectoryUncompressed contains the uncompressed build context
	ContextDirectoryUncompressed = "/beval/context/uncompressed/"
	// Remote Local Directory
	RemoteLocalDirectory = "."
	// Remote Temp Directory
	RemoteTempDirectory = "/ocib/temp/"
)

// Remote paths
const (
	OverlayPath    = "./overlay_DOWNLOAD.yaml"
	DockerStepPath = "./step_cmds_DOWNLOAD.yaml"
)
