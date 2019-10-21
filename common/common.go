package common

import "github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder"

// Controller environment variables
const (
	// EnvVarControllerConfigMap is the name of the configmap to use for the controller
	EnvVarControllerConfigMap = "CONTROLLER_CONFIG_MAP"
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
	// LabelNamespace is the label to indicate K8s namespace
	LabelNamespace = "namespace"
)

// Miscellaneous constants for controller
const (
	// ControllerConfigMapKey is the key in the configmap to retrieve ocibuilder controller configuration from.
	// Content encoding is expected to be YAML.
	ControllerConfigMapKey = "config"
)

// logger labels
const (
	LabelVersion = "version"
)
