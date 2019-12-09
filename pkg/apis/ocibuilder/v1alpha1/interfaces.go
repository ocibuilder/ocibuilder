package v1alpha1

import (
	"io"

	"github.com/docker/docker/api/types"
)

// Overlay is the overlay interface for handling of overylaying onto specification files
type Overlay interface {
	// Apply applies an overlay in your implementing struct
	Apply() ([]byte, error)
}

// ContextReader provides an interface for reading from multiple different build contexts
type ContextReader interface {
	Read() (io.ReadCloser, error)
}

// SpecGenerator provides an interface for spec generation for ocibuilder.yaml specification files
type SpecGenerator interface {
	Generate() ([]byte, error)
}

type BuilderClient interface {
	ImageBuild(options OCIBuildOptions) (OCIBuildResponse, error)
	ImagePull(options OCIPullOptions) (OCIPullResponse, error)
	ImagePush(options OCIPushOptions) (OCIPushResponse, error)
	ImageRemove(options OCIRemoveOptions) (OCIRemoveResponse, error)
	RegistryLogin(options OCILoginOptions) (OCILoginResponse, error)
	GenerateAuthRegistryString(auth types.AuthConfig) string
}
