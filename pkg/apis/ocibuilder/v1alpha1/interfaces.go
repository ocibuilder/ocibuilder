package v1alpha1

import "io"

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
