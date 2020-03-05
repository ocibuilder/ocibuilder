/*
Copyright © 2019 BlackRock Inc.

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
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/ocibuilder/ocibuilder/pkg/command"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodePhase is the label for the condition of a node.
type NodePhase string

// possible types of node phases
const (
	NodePhaseRunning   NodePhase = "Running"   // the node is running
	NodePhaseError     NodePhase = "Error"     // the node has encountered an error in processing
	NodePhaseNew       NodePhase = ""          // the node is new
	NodePhaseCompleted NodePhase = "Completed" // node has completed running
)

// Framework is the type of the build framework being used
type Framework string

const (
	// DockerFramework is the type of docker framework
	DockerFramework Framework = "docker"
	// BuildahFramework is the type of buildah framework
	BuildahFramework Framework = "buildah"
)

const (
	// AnsibleTemplateDir is the path for ansible template
	AnsibleTemplateDir string = "../../templates/ansible"
	// AnsibleTemplate is the path for ansible template
	AnsibleTemplate string = "ansible.tmpl"
	// AnsibleBase is the ansible base directory
	AnsibleBase string = "/etc/ansible"
)

// MetadataType is the type of metadata that you want to store
type MetadataType string

const (
	// Build is the build related metadata type
	Build MetadataType = "build"
	// Attestation is attestation metadata
	Attestation MetadataType = "attestation"
	// DerviedImage any metadata related to the derived image
	Image MetadataType = "image"
)

// OCIBuilder is the definition of a ocibuilder resource
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type OCIBuilder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline" protobuf:"bytes,1,name=metadata"`
	Spec              OCIBuilderSpec   `json:"spec" protobuf:"bytes,2,name=spec"`
	Status            OCIBuilderStatus `json:"status" protobuf:"bytes,3,name=status"`
}

// OCIBuilderList is the list of OCIBuilder resources.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OCIBuilderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,name=metadata"`
	// +listType=map
	Items []OCIBuilder `json:"items" protobuf:"bytes,2,name=items"`
}

// OCIBuilderSpec represents OCIBuilder specifications.
type OCIBuilderSpec struct {
	// Envs are the list of environment variables available to components.
	// +optional
	// +listType=map
	Params []Param `json:"params,omitempty" protobuf:"bytes,1,opt,name=params"`
	// Logins holds information to log into one or more registries
	// +listType=map
	Login []LoginSpec `json:"login,omitempty" protobuf:"bytes,2,opt,name=login"`
	// Build represents the build specifications for images
	// +optional
	Build *BuildSpec `json:"build,omitempty" protobuf:"bytes,3,name=build"`
	// Push contains specification to push images to registries
	// +optional
	// +listType=map
	Push []PushSpec `json:"push,omitempty" protobuf:"bytes,4,name=push"`
	// Type of the build framework.
	// Defaults to docker
	// +optional
	Daemon bool `json:"daemon,omitempty" protobuf:"bytes,5,opt,name=daemon"`
	// Configuration for storing build metadata in an external Metadata store.
	// Defaults to Grafeas as the chosen metadata store
	// +optional
	Metadata *Metadata `json:"metadata,omitempty" protobuf:"bytes,6,opt,name=metadata"`
}

// OCIBuilderStatus holds the status of a OCIBuilder resource
type OCIBuilderStatus struct {
	// Phase is the high-level summary of the OCIBuilder
	Phase NodePhase `json:"phase" protobuf:"bytes,1,opt,name=phase"`
	// StartedAt is the time at which this OCIBuilder was initiated
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,2,opt,name=startedAt"`
	// Message is a human readable string indicating details about a OCIBuilder in its phase
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Nodes is a mapping between a node ID and the node's status
	// it records the states for the configurations of OCIBuilder.
	Nodes map[string]*NodeStatus `json:"nodes" protobuf:"bytes,1,name=nodes"`
}

// Param represents parameters
type Param struct {
	// Value of the environment variable.
	// +optional
	Value string `json:"value,omitempty" protobuf:"bytes,1,opt,name=value"`
	// Dest is the destination of the field to replace with the parameter
	Dest string `json:"dest" protobuf:"bytes,2,opt,name=dest"`
	// ValueFromEnvVar is a variable which is to be replaced by an env var
	// +optional
	ValueFromEnvVariable string `json:"valueFromEnvVariable,omitempty" protobuf:"bytes,3,opt,name=valueFromEnvVariable"`
}

// BuildSpec represents the build specifications for images
type BuildSpec struct {
	// Templates are set of build templates that describe steps needed to build a Dockerfile
	// +listType=map
	Templates []BuildTemplate `json:"templates" protobuf:"bytes,1,rep,name=templates"`
	// Steps within a build
	// +listType=map
	Steps []BuildStep `json:"steps" protobuf:"bytes,2,rep,name=steps"`
	// StorageDriver is the storage driver flag (default overlay2) see https://docs.docker.com/storage/storagedriver/select-storage-driver/
	StorageDriver string `json:"storageDriver" protobuf:"bytes,2,rep,name=storageDriver"`
}

// BuildTemplate represents the build template that can shared across different builds
type BuildTemplate struct {
	// Name of the template
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// List of cmds in a Dockerfile
	// +listType=
	Cmd []BuildTemplateStep `json:"cmd" protobuf:"bytes,2,rep,name=steps"`
}

// BuildTemplateStep represents a step within build template
type BuildTemplateStep struct {
	// Docker represents a docker step within build template steps
	Docker *DockerStep `json:"docker,omitempty" protobuf:"bytes,1,opt,name=docker"`
	// Ansible represents a ansible step within build template steps
	Ansible *AnsibleStep `json:"ansible,omitempty" protobuf:"bytes,2,opt,name=ansible"`
}

// DockerStep represents a step within a build that contains docker commands
type DockerStep struct {
	// Inline Dockerfile commands
	// +optional
	// +listType=map
	Inline []string `json:"inline,omitempty" protobuf:"bytes,1,opt,name=inline"`
	// Path to a file that contains Dockerfile commands
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`
	// Remote url to a file that contains docker commands
	// +optional
	Url string `json:"url,omitempty" protobuf:"bytes,3,opt,name=url"`
	// Auth for remote access to a url
	// +optional
	Auth RemoteCreds `json:"auth,inline" protobuf:"bytes,4,name=auth"`
}

// AnsibleStep represents an ansible install  within a build
type AnsibleStep struct {
	// Playbook refers to playbook.yaml file
	Playbook string `json:"playbook" protobuf:"bytes,1,name=playbook"`
	// Requirements refer to the requirements.yaml file
	// +optional
	Requirements string `json:"requirements,omitempty" protobuf:"bytes,2,opt,name=requirements"`
	// Workspace is the name of your ansible workspce NOT including /etc/ansible/ ansible path
	Workspace string `json:"workspace" protobuf:"bytes,3,name=workspace"`
}

// BuildStep represents a step within the build
type BuildStep struct {
	// Metadata about the build step.
	*ImageMetadata `json:"metadata,inline" protobuf:"bytes,1,name=metadata"`
	// Stages of the build
	// +listType=map
	Stages []Stage `json:"stages" protobuf:"bytes,3,opt,name=purge"`
	// Git url to fetch the project from.
	// +optional
	Git string `json:"git,omitempty" protobuf:"bytes,3,opt,name=git"`
	// Tag the tag of the build
	// +optional
	Tag string `json:"tag,omitempty" protobuf:"bytes,4,opt,name=tag"`
	// Distroless if set to true generates a distroless image
	Distroless bool `json:"distroless,omitempty" protobuf:"bytes,5,opt,name=distroless"`
	// Cache for build
	// Set to false by default
	// +optional
	Cache bool `json:"cache,omitempty" protobuf:"bytes,6,opt,name=cache"`
	// Purge the build
	// defaults to false
	// +optional
	Purge bool `json:"purge,omitempty" protobuf:"bytes,7,opt,name=purge"`
	// BuildContext used for image build
	// default looks at the current working directory
	// +optional
	BuildContext *BuildContext `json:"context,omitempty" protobuf:"bytes,8,opt,name=context"`
}

// Stage represents a stage within the build
type Stage struct {
	// Metadata refers to metadata of the build stage
	*ImageMetadata `json:"metadata,inline" protobuf:"bytes,1,name=metadata"`
	// BaseImage refers to parent image for given build stage.
	Base Base `json:"base" protobuf:"bytes,2,name=base"`
	// Template refers to one of the build templates.
	Template string `json:"template" protobuf:"bytes,3,name=template"`
	// Cmd refers to a template defined in a stage without a template.
	// +listType=map
	Cmd []BuildTemplateStep `json:"cmd" protobuf:"bytes,4,name=cmd"`
}

// Base represents base image details
type Base struct {
	Image string `json:"image" protobuf:"bytes,1,name=image"`
	// Tag is the tag for the image
	// +optional
	Tag string `json:"tag,omitempty" protobuf:"bytes,2,name=tag"`
	// Platform is the specified platform of the image
	// +optional
	Platform string `json:"platform,omitempty" protobuf:"bytes,3,name=platform"`
}

// ImageMetadata represents data about a build step
type ImageMetadata struct {
	// Name of the build step
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// Labels for the step
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,2,opt,name=labels"`
	// Annotations for the step
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,3,opt,name=annotations"`
	// Creator is the creator of the build
	Creator string `json:"creator,omitempty" protobuf:"bytes,4,opt,name=creator"`
	// Source is the URI to the source code of the image build
	Source string `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`
}

// LoginSpec holds the information to log into a registry.
type LoginSpec struct {
	// Registry refers to a OCI image registry
	Registry string `json:"registry" protobuf:"bytes,1,name=registry"`
	Token    string `json:"token" protobuf:"bytes,2,name=token"`
	// Creds refer to credentials required to log into the registry
	Creds RegistryCreds `json:"creds" protobuf:"bytes,3,name=creds"`
	// Overlay is the name which will be referred to by an overlay file
	Overlay string `json:"overlay" protobuf:"bytes,4,name=overlay"`
}

// RegistryCreds holds the credentials to login into a registry
type RegistryCreds struct {
	// K8s refer to the credentials stored in K8s secrets
	K8s K8sCreds `json:"k8s,omitempty" protobuf:"bytes,1,opt,name=k8s"`
	// Env refers to the credentials stored in environment variables
	Env EnvCreds `json:"env,omitempty" protobuf:"bytes,2,opt,name=env"`
	// Plain refers to the credentials set inline
	Plain PlainCreds `json:"plain,omitempty" protobuf:"bytes,3,opt,name=plain"`
}

// RemoteCreds holds the credentials to pull from a remote url
type RemoteCreds struct {
	// Env refers to the credentials stored in environment variables
	Env EnvCreds `json:"env,omitempty" protobuf:"bytes,2,opt,name=env"`
	// Plain refers to the credentials set inline
	Plain PlainCreds `json:"plain,omitempty" protobuf:"bytes,3,opt,name=plain"`
}

// K8sCreds refers to the K8s secret that holds the registry creds.
type K8sCreds struct {
	// Username refers to the K8s secret that holds username
	Username *corev1.SecretKeySelector `json:"username" protobuf:"bytes,1,name=username"`
	// Password refers to the K8s secret that holds password
	Password *corev1.SecretKeySelector `json:"password" protobuf:"bytes,2,name=password"`
}

// EnvCreds refers to credentials stored in env vars.
type EnvCreds struct {
	// Username refers to an env var that holds the username
	Username string `json:"username" protobuf:"bytes,1,name=username"`
	// Password refers to an en var that holds the password
	Password string `json:"password" protobuf:"bytes,2,name=password"`
}

// PlainCreds refers to the credentials set inline
type PlainCreds struct {
	Username string `json:"username" protobuf:"bytes,1,name=username"`
	Password string `json:"password"`
}

// PushSpec contains the specification to push images to registries
type PushSpec struct {
	// Registry is the name of the registry
	Registry string `json:"registry" protobuf:"bytes,1,name=registry"`
	// Image to push
	Image string `json:"image" protobuf:"bytes,2,name=image"`
	// User is the name of kubernetes namespace
	User string `json:"user" protobuf:"bytes,3,name=user"`
	// Token required for the OCI complaint registry authentication
	Token string `json:"token" protobuf:"bytes,4,name=token"`
	// Tag version of the image (e.g: v0.1.1)
	Tag string `json:"tag" protobuf:"bytes,5,name=tag"`
	// Purge the image after it has been pushed
	// defaults to false
	// +optional
	Purge bool `json:"purge,omitempty" protobuf:"bytes,6,opt,name=purge"`
	// Overlay is the name which will be referred to by an overlay file
	Overlay string `json:"overlay" protobuf:"bytes,7,name=overlay"`
}

// NodeStatus describes the status for an individual node in the ocibuilder configurations.
// A single node can represent one configuration.
type NodeStatus struct {
	// ID is a unique identifier of a node within build steps
	// It is a hash of the node name
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`
	// Name is a unique name in the node tree used to generate the node ID
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// DisplayName is the human readable representation of the node
	DisplayName string `json:"displayName" protobuf:"bytes,5,opt,name=displayName"`
	// Phase of the node
	Phase NodePhase `json:"phase" protobuf:"bytes,6,opt,name=phase"`
	// StartedAt is the time at which this node started
	StartedAt metav1.MicroTime `json:"startedAt,omitempty" protobuf:"bytes,7,opt,name=startedAt"`
	// Message store data or something to save for configuration
	Message string `json:"message,omitempty" protobuf:"bytes,8,opt,name=message"`
	// UpdateTime is the time when node(OCIBuilder configuration) was updated
	UpdateTime metav1.MicroTime `json:"updateTime,omitempty" protobuf:"bytes,9,opt,name=updateTime"`
}

// ImageBuildArgs describes the arguments for running a build command
type ImageBuildArgs struct {
	// Name is the name of the build
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Tag is the tag of the build
	Tag string `json:"tag" protobuf:"bytes,2,opt,name=tag"`
	// Dockerfile is the path to the generated Dockerfile
	// +optional
	Dockerfile string `json:"dockerfile,omitempty" protobuf:"bytes,3,opt,name=dockerfile"`
	// Purge the image after it has been pushed
	// defaults to false
	// +optional
	Purge bool `json:"purge,omitempty" protobuf:"bytes,5,opt,name=purge"`
	// BuildContextPath is the path of the build context for Docker and Buildah
	// defaults to LocalContext in current working directory
	// +optional
	BuildContextPath string `json:"buildContextPath,omitempty" protobuf:"bytes,6,opt,name=buildContextPath"`
	// Labels for the step
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,7,opt,name=labels"`
	// Creator is the email of the build creator
	Creator string `json:"creator,omitempty"`
	// Source is the URI of the source code for the build
	Source string `json:"source,omitempty"`
	// Cache for build
	// Set to false by default
	// +optional
	Cache bool `json:"cache,omitempty" protobuf:"bytes,8,opt,name=cache"`
	// StorageDriver is a buildah flag for storage driver e.g. vfs
	StorageDriver string `json:"storageDriver" protobuf:"bytes,9,name=storageDriver"`
}

// BuildContext stores the chosen build context for your build, this can be Local, S3 or Git
type BuildContext struct {
	// Local context contains local context information for a build
	LocalContext *LocalContext `json:"localContext,omitempty" protobuf:"bytes,1,opt,name=localContext"`
	// S3Context refers to the context stored on S3 bucket for a build
	S3Context *S3Context `json:"s3Context,omitempty" protobuf:"bytes,2,opt,name=s3Context"`
	// GitContext refers to the context stored on Git repository
	GitContext *GitContext `json:"gitContext,omitempty" protobuf:"bytes,3,opt,name=gitContext"`
	// GCSContext refers to the context stored on the GCS
	GCSContext *GCSContext `json:"gcsContext,omitempty" protobuf:"bytes,4,opt,name=gcsContext"`
	// AzureBlobContext refers to the context stored on the Azure Storage Blob
	AzureBlobContext *AzureBlobContext `json:"azureBlobContext,omitempty" protobuf:"bytes,5,opt,name=azureBlobContext"`
	// AliyunOSSContext refers to the context stored on the Aliyun OSS
	AliyunOSSContext *AliyunOSSContext `json:"aliyunOSSContext,omitempty" protobuf:"bytes,6,opt,name=aliyunOSSContext"`
}

// LocalContext stores the path for your local build context, implements the ContextReader interface
type LocalContext struct {
	// ContextPath is the path to your build context
	ContextPath string `json:"contextPath" protobuf:"bytes,1,opt,name=contextPath"`
}

// KubeSecretCredentials refers to K8s secret that holds the credentials
type KubeSecretCredentials struct {
	// Secret is the K8s secret key selector
	Secret *corev1.SecretKeySelector `json:"secret" protobuf:"bytes,1,name=secret"`
	// Namespace where the secret is stored
	Namespace string `json:"namespace" protobuf:"bytes,2,name=namespace"`
}

// Credentials encapsulates different ways of storing the credentials
type Credentials struct {
	// Plain text credentials
	Plain string `json:"plain,omitempty" protobuf:"bytes,1,opt,name=plain"`
	// Env refers to credentials stored in environment variable
	Env string `json:"env,omitempty" protobuf:"bytes,2,opt,name=env"`
	// KubeSecret refers to K8s secret that holds the credentials
	KubeSecret *KubeSecretCredentials `json:"kubeSecret,omitempty" protobuf:"bytes,3,opt,name=kubeSecret"`
}

// S3Bucket contains information to describe an S3 Bucket
type S3Bucket struct {
	Key  string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

// S3Context refers to context stored on S3 bucket to build an image
type S3Context struct {
	Endpoint  string       `json:"endpoint" protobuf:"bytes,1,name=endpoint"`
	Bucket    *S3Bucket    `json:"bucket" protobuf:"bytes,2,name=bucket"`
	Region    string       `json:"region,omitempty" protobuf:"bytes,3,opt,name=region"`
	Insecure  bool         `json:"insecure,omitempty" protobuf:"variant,4,opt,name=insecure"`
	AccessKey *Credentials `json:"accessKey" protobuf:"bytes,5,name=accessKey"`
	SecretKey *Credentials `json:"secretKey" protobuf:"bytes,6,name=secretKey"`
}

// GitRemoteConfig contains the configuration of a Git remote
type GitRemoteConfig struct {
	// Name of the remote to fetch from.
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// +listType=urls
	// URLs the URLs of a remote repository. It must be non-empty. Fetch will
	// always use the first URL, while push will use all of them.
	URLS []string `json:"urls" protobuf:"bytes,2,rep,name=urls"`
}

// GitContext contains information about an artifact stored in git
type GitContext struct {
	// Git URL
	URL string `json:"url" protobuf:"bytes,1,name=url"`
	// Username for authentication
	Username *Credentials `json:"username,omitempty" protobuf:"bytes,2,opt,name=username"`
	// Password for authentication
	Password *Credentials `json:"password,omitempty" protobuf:"bytes,3,opt,name=password"`
	// SSHKeyPath is path to your ssh key path. Use this if you don't want to provide username and password.
	// ssh key path must be mounted in sensor pod.
	// +optional
	SSHKeyPath string `json:"sshKeyPath,omitempty" protobuf:"bytes,4,opt,name=sshKeyPath"`
	// Branch to use to pull trigger resource
	// +optional
	Branch string `json:"branch,omitempty" protobuf:"bytes,5,opt,name=branch"`
	// Tag to use to pull trigger resource
	// +optional
	Tag string `json:"tag,omitempty" protobuf:"bytes,6,opt,name=tag"`
	// Ref to use to pull trigger resource. Will result in a shallow clone and
	// fetch.
	// +optional
	Ref string `json:"ref,omitempty" protobuf:"bytes,7,opt,name=ref"`
	// Remote to manage set of tracked repositories. Defaults to "origin".
	// Refer https://git-scm.com/docs/git-remote
	// +optional
	Remote *GitRemoteConfig `json:"remote" protobuf:"bytes,8,opt,name=remote"`
}

// GCSContext refers to the context stored on GCP Storage
type GCSContext struct {
	// CredentialsFilePath refers to the credentials file path
	CredentialsFilePath string `json:"credentialsFilePath,omitempty" protobuf:"bytes,1,opt,name=credentialsFilePath"`
	// APIKey for authentication
	APIKey *Credentials `json:"apiKey,omitempty" protobuf:"bytes,2,opt,name=apiKey"`
	// AuthRequired checks if authentication is required to connect to GCS
	AuthRequired bool `json:"authRequired" protobuf:"bytes,3,name=authRequired"`
	// Endpoint is the storage to connect to
	Endpoint string `json:"endpoint" protobuf:"bytes,4,name=endpoint"`
	// Bucket refers to the bucket name on gcs
	Bucket *S3Bucket `json:"bucket" protobuf:"bytes,5,name=bucket"`
	// Region refers to GCS region
	Region string `json:"region,omitempty" protobuf:"bytes,6,opt,name=region"`
}

// AzureBlobContext refers to configuration required to fetch context from Azure Storage Blob
type AzureBlobContext struct {
	// AzureStorageAccount refers to the account name
	Account *Credentials `json:"account" protobuf:"bytes,1,name=account"`
	// AccessKey refers to the access key
	AccessKey *Credentials `json:"accessKey" protobuf:"bytes,2,name=accessKey"`
	// URL refers to blob's URL
	URL *Credentials `json:"url" protobuf:"bytes,3,name=url"`
}

// AliyunOSSContext refers to configuration required to fetch context from Aliyun OSS
type AliyunOSSContext struct {
	// AccessId refers to access id
	AccessId *Credentials `json:"accessId" protobuf:"bytes,1,name=accessId"`
	// AccessSecret refers to access secret
	AccessSecret *Credentials `json:"accessSecret" protobuf:"bytes,2,name=accessSecret"`
	// Endpoint is the storage to connect to
	Endpoint string `json:"endpoint" protobuf:"bytes,4,name=endpoint"`
	// Bucket refers to the bucket name on gcs
	Bucket *S3Bucket `json:"bucket" protobuf:"bytes,5,name=bucket"`
}

// Command Represents a single line in a Dockerfile
type Command struct {
	// Cmd lowercased command name (e.g `from`)
	Cmd string `json:"cmd" protobuf:"bytes,1,opt,name=cmd"`
	// SubCmd for ONBUILD only this holds the sub-command
	SubCmd string `json:"subCmd" protobuf:"bytes,2,opt,name=subCmd"`
	// Json bool for whether the value is written in json
	IsJSON bool `json:"isJSON" protobuf:"bytes,3,opt,name=isJSON"`
	// Original is the original source line
	Original string `json:"original" protobuf:"bytes,4,opt,name=original"`
	// StartLine is the original source line number
	StartLine int `json:"startLine" protobuf:"bytes,5,opt,name=startLine"`
	// Flags such as `--from=...` for `COPY`.
	// +listType=map
	Flags []string `json:"flags" protobuf:"bytes,6,opt,name=flags"`
	// Value is the contents of the command (e.g `ubuntu:xenial`)
	// +listType=map
	Value []string `json:"value" protobuf:"bytes,7,opt,name=value"`
}

// BuildProvenance represents build image metadata
type BuildProvenance struct {
	// BuildFile is the path to the buildfile that was used for the image build
	BuildFile string `json:"buildFile" protobuf:"bytes,1,opt,name=buildFile"`
	// ContextDirectory is the path to the build context
	ContextDirectory string `json:"contextDirectory" protobuf:"bytes,2,opt,name=contextDirectory"`
	// Daemon is whether the daemon was used to build or not (Docker or Buildah)
	Daemon bool `json:"daemon" protobuf:"bytes,3,opt,name=daemon"`
	// Time at which the build was created.
	CreateTime time.Time `json:"createTime,omitempty"`
	// Time at which execution of the build was started.
	StartTime time.Time `json:"startTime,omitempty"`
	// Time at which execution of the build was finished.
	EndTime time.Time `json:"endTime,omitempty"`
	// Creator is the email of the build creator
	Creator string `json:"creator,omitempty"`
	// Source is the URI of the source code for the build
	Source string `json:"source,omitempty"`
	// Name is the image name
	Name string `json:"name,omitempty"`
	// Tag is the image tag
	Tag string `json:"tag,omitempty"`
	// ID is the ID of the image
	ID string `json:"id,omitempty"`
}

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
	// Finished is the flag to determine that the response has finished being read
	Finished bool
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
	// Finished is the flag to determine that the response has finished being read
	Finished bool
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

// GenerateTemplate is the template for a docker generate
type GenerateTemplate struct {
	ImageName string
	Tag       string
	Stages    []string
	Templates []string
}

// StageGenTemplate is the template for a stage in docker generate
type StageGenTemplate struct {
	Base         string
	BaseTag      string
	StageName    string
	TemplateName string
}

// BuildGenTemplate is the template for a build template in docker generate
type BuildGenTemplate struct {
	Name string
	Cmds []string
}

// Metadata is where metadata to store is defined in the ocibuilder specification
type Metadata struct {
	// StoreType is the metadata store type to push metadata to
	StoreConfig *StoreConfig `json:"storeConfig,omitempty" protobuf:"bytes,1,opt,name=storeConfig"`
	// SignKey holds the key to sign an image for attestation purposes
	Key *SignKey `json:"signKey,omitempty" protobuf:"bytes,2,opt,name=signKey"`
	// Hostname is the hostname of the metadatastore
	Hostname string `json:"hostname,omitempty" protobuf:"bytes,3,opt,name=hostname"`
	// Data is the types of metadata that you would like to push to your metadatastore
	Data []MetadataType `json:"data,omitempty" protobuf:"bytes,4,opt,name=data"`
	// Creator is the email of the build creator
	Creator string `json:"creator,omitempty" protobuf:"bytes,5,opt,name=creator"`
}

type SignKey struct {
	// PrivateKey is an ascii armored private key used to sign images for image attestation
	// +optional
	PlainPrivateKey string `json:"plainPrivateKey,omitempty" protobuf:"bytes,1,opt,name=plainPrivateKey"`
	// PublicKey is the ascii armored public key for verification in image attestation
	// +optional
	// +optional
	PlainPublicKey string `json:"plainPublicKey,omitempty" protobuf:"bytes,2,opt,name=plainPublicKey"`
	// EnvPrivateKey is an env variable that holds an ascii armored private key used to sign images for image attestation
	// +optional
	EnvPrivateKey string `json:"envPrivateKey,omitempty" protobuf:"bytes,3,opt,name=envPrivateKey"`
	// EnvPublicKey is an env variable that holds an ascii armored public key used to sign images for image attestation
	// +optional
	EnvPublicKey string `json:"envPublicKey,omitempty" protobuf:"bytes,4,opt,name=envPublicKey"`
	// Passphrase is the passphrase for decrypting the private key
	Passphrase string `json:"passphrase,omitempty" protobuf:"bytes,5,opt,name=passphrase"`
	// Url or a filepath to a file that contains an ascii armored private key
	// +optional
	Url string `json:"url,omitempty" protobuf:"bytes,6,opt,name=url"`
	// Auth for remote access to a url
	// +optional
	Auth RemoteCreds `json:"auth,inline" protobuf:"bytes,7,name=auth"`
}

// StoreConfig is the configuration of the metadata store to push metadata to
type StoreConfig struct {
	// Grafeas holds the config for the Grafeas metadata store
	Grafeas *Grafeas `json:"grafeas,omitempty" protobuf:"bytes,1,opt,name=grafeas"`
}

// Grafeas is the type defining the Grafeas metadata store
type Grafeas struct {
	// Project is the name of the project ID to store the occurrence
	Project string `json:"project,omitempty" protobuf:"bytes,1,opt,name=project"`
	// Notes holds the notes for the three occurrence types
	Notes Notes `json:"notes,omitempty" protobuf:"bytes,3,opt,name=notes"`
}

type Notes struct {
	// BuildNoteName Required. Immutable. The analysis note associated with build occurrence, in the form of `projects/[PROVIDER_ID]/notes/[NOTE_ID]`. This field can be used as a filter in list requests.
	BuildNoteName string `json:"build,omitempty" protobuf:"bytes,1,opt,name=build"`
	// AttestationNoteName Required. Immutable. The analysis note associated with attestation occurrence, in the form of `projects/[PROVIDER_ID]/notes/[NOTE_ID]`. This field can be used as a filter in list requests.
	AttestationNoteName string `json:"attestation,omitempty" protobuf:"bytes,2,opt,name=attestation"`
	// DerivedImageNoteName Required. Immutable. The analysis note associated with image derived occurrence, in the form of `projects/[PROVIDER_ID]/notes/[NOTE_ID]`. This field can be used as a filter in list requests.
	DerivedImageNoteName string `json:"image,omitempty" protobuf:"bytes,3,opt,name=image"`
}
