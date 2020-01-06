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

package v1alpha1

import (
	"github.com/docker/docker/api/types"
	"io"
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

// BuilderClient is the client interface for the ocibuilder
type BuilderClient interface {
	ImageBuild(options OCIBuildOptions) (OCIBuildResponse, error)
	ImagePull(options OCIPullOptions) (OCIPullResponse, error)
	ImagePush(options OCIPushOptions) (OCIPushResponse, error)
	ImageRemove(options OCIRemoveOptions) (OCIRemoveResponse, error)
	RegistryLogin(options OCILoginOptions) (OCILoginResponse, error)
	GenerateAuthRegistryString(auth types.AuthConfig) string
}
