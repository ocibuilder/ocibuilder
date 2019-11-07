package v1alpha1

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
)

// Builder provides an interface for both Buildah and Docker build tools
type Builder interface {
	Build(spec OCIBuilderSpec) ([]io.ReadCloser, error)
	Login(spec OCIBuilderSpec) ([]io.ReadCloser, error)
	Pull(spec OCIBuilderSpec, name string) ([]io.ReadCloser, error)
	Push(spec OCIBuilderSpec) ([]io.ReadCloser, error)
	Clean()
}

// Overlay is the overlay interface for handling of overylaying onto specification files
type Overlay interface {
	// Apply applies an overlay in your implementing struct
	Apply() ([]byte, error)
}

// ContextReader provides an interface for reading from multiple different build contexts
type ContextReader interface {
	Read() (io.ReadCloser, error)
}

type BuilderClient interface {
	ImageBuild(options OCIBuildOptions) (types.ImageBuildResponse, error)
	ImagePull(options OCIPullOptions) (io.ReadCloser, error)
	ImagePush(options OCIPushOptions) (io.ReadCloser, error)
	ImageRemove(options OCIRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	RegistryLogin(options OCILoginOptions) (registry.AuthenticateOKBody, error)
	GenerateAuthRegistryString(registry string) (string, error)
}
