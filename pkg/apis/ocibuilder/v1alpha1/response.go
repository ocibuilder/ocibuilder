/*
Copyright © 2020 BlackRock Inc.

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
	ctx "context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/ocibuilder/ocibuilder/pkg/command"
)

// OCIBuildOptions are the build options for an ocibuilder build
type OCIBuildOptions struct {
	// ImageBuildOptions are standard Docker API image build options
	types.ImageBuildOptions `json:"imageBuildOptions,inline" protobuf:"bytes,1,name=imageBuildOptions"`
	// ContextPath is the path to the raw build context, used for Buildah builds
	ContextPath string `json:"contextPath" protobuf:"bytes,2,name=contextPath"`
	// Ctx is the goroutine context
	Ctx ctx.Context `json:"ctx" protobuf:"bytes3,name=ctx"`
	// Context is the docker tared build context
	Context io.Reader `json:"context" protobuf:"bytes,4,name=context"`
	// StorageDriver is a buildah flag for storage driver e.g. vfs
	StorageDriver string `json:"storageDriver" protobuf:"bytes,5,name=storageDriver"`
}

// OCIBuildResponse is the build response from an ocibuilder build
type OCIBuildResponse struct {
	// ImageBuildResponse is standard build response from the Docker API
	types.ImageBuildResponse `json:"imageBuildResponse,inline" protobuf:"bytes,1,name=imageBuildResponse"`
	// Exec is part of the response for Buildah command executions
	Exec *command.Command `json:"exec,inline" protobuf:"bytes,2,name=exec"`
	// Stderr is the stderr output stream used to stream buildah response
	Stderr io.ReadCloser `json:"stderr,inline" protobuf:"bytes,3,name=stderr"`
}

// OCIPullOptions are the pull options for an ocibuilder pull
type OCIPullOptions struct {
	// ImagePullOptions are the standard Docker API pull options
	types.ImagePullOptions `json:"imagePullOptions,inline" protobuf:"bytes,1,name=imagePullOptions"`
	// Ref is the reference image name to pull
	Ref string `json:"ref,inline" protobuf:"bytes,2,name=ref"`
	// Ctx is the goroutine context
	Ctx ctx.Context `json:"ctx,inline" protobuf:"bytes3,name=ctx"`
}

// OCIPullResponse is the pull response from an ocibuilder pull
type OCIPullResponse struct {
	// Body is the body of the response from an ocibuilder pull
	Body io.ReadCloser `json:"body,inline" protobuf:"bytes,1,name=body"`
	// Exec is part of the response for Buildah command executions
	Exec *command.Command `json:"exec,inline" protobuf:"bytes,2,name=exec"`
	// Stderr is the stderr output stream used to stream buildah response
	Stderr io.ReadCloser `json:"stderr,inline" protobuf:"bytes,3,name=stderr"`
}

// OCIPushOptions are the pull options for an ocibuilder push
type OCIPushOptions struct {
	// ImagePushOptions are the standard Docker API push options
	types.ImagePushOptions `json:"imagePushOptions,inline" protobuf:"bytes,1,name=imagePushOptions"`
	// Ref is the reference image name to push
	Ref string `json:"ref,inline" protobuf:"bytes,2,name=ref"`
	// Ctx is the goroutine context
	Ctx ctx.Context `json:"ctx,inline" protobuf:"bytes3,name=ctx"`
}

// OCIPushResponse is the push response from an ocibuilder push
type OCIPushResponse struct {
	// Body is the body of the response from an ocibuilder push
	Body io.ReadCloser `json:"body,inline" protobuf:"bytes,1,name=body"`
	// Exec is part of the response for Buildah command executions
	Exec *command.Command `json:"exec,inline" protobuf:"bytes,2,name=exec"`
	// Stderr is the stderr output stream used to stream buildah response
	Stderr io.ReadCloser `json:"stderr,inline" protobuf:"bytes,3,name=stderr"`
}

// OCIRemoveOptions are the remove options for an ocibuilder remove
type OCIRemoveOptions struct {
	// ImageRemoveOptions are the standard Docker API remove options
	types.ImageRemoveOptions `json:"imageRemoveOptions,inline" protobuf:"bytes,1,name=imageRemoveOptions"`
	// Image is the name of the image to remove
	Image string `json:"image,inline" protobuf:"bytes,2,name=image"`
	// Ctx is the goroutine context
	Ctx ctx.Context `json:"ctx,inline" protobuf:"bytes,3,name=ctx"`
}

// OCIRemoveResponse is the response from an ocibuilder remove
type OCIRemoveResponse struct {
	// Response are the responses from an image delete
	Response []types.ImageDeleteResponseItem `json:"response,inline" protobuf:"bytes,1,name=response"`
	// Exec is part of the response for Buildah command executions
	Exec *command.Command `json:"exec,inline" protobuf:"bytes,2,name=exec"`
	// Stderr is the stderr output stream used to stream buildah response
	Stderr io.ReadCloser `json:"stderr,inline" protobuf:"bytes,3,name=stderr"`
}

// OCILoginOptions are the login options for an ocibuilder login
type OCILoginOptions struct {
	// AuthConfig is the standard auth config for the Docker API
	types.AuthConfig `json:"authConfig,inline" protobuf:"bytes,1,name=authConfig"`
	// Ctx is the goroutine context
	Ctx ctx.Context `json:"ctx,inline" protobuf:"bytes,2,name=ctx"`
}

// OCILoginResponse is the login response from an ocibuilder login
type OCILoginResponse struct {
	// AuthenticateOKBody is the standar login response from the Docker API
	registry.AuthenticateOKBody
	// Exec is part of the response for Buildah command executions
	Exec *command.Command `json:"exec,inline" protobuf:"bytes,2,name=exec"`
	// Stderr is the stderr output stream used to stream buildah response
	Stderr io.ReadCloser `json:"stderr,inline" protobuf:"bytes,3,name=stderr"`
}
