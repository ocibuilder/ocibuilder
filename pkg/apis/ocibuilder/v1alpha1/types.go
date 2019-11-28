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
	// AnsiblePath is the path for ansible module
	AnsiblePath string = "ansible"
	// AnsibleGalaxyPath is the path of ansible galaxy
	AnsibleGalaxyPath string = "ansible-galaxy"
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
}

// AnsibleStep represents an ansible install  within a build
type AnsibleStep struct {
	// Local contains information to install a ansible role through local playbook
	Local *AnsibleLocal `json:"local,omitempty" protobuf:"bytes,1,opt,name=local"`
	// Galaxy contains information to install a ansible role through ansible-galaxy
	Galaxy *AnsibleGalaxy `json:"galaxy,omitempty" protobuf:"bytes,2,opt,name=galaxy"`
}

// AnsibleLocal contains information to install a ansible role through local playbook
type AnsibleLocal struct {
	// Playbook refers to playbook.yaml file
	Playbook string `json:"playbook" protobuf:"bytes,2,name=playbook"`
}

// AnsibleGalaxy contains the information about the role to install through galaxy
type AnsibleGalaxy struct {
	// Requirements refer to the requirements.yaml file
	// +optional
	Requirements string `json:"requirements,omitempty" protobuf:"bytes,1,opt,name=requirements"`
	// Name of the galaxy role
	Name string `json:"name" protobuf:"bytes,2,name=name"`
}

// BuildStep represents a step within the build
type BuildStep struct {
	// Metadata about the build step.
	*Metadata `json:"metadata,inline" protobuf:"bytes,1,name=metadata"`
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
	*Metadata `json:"metadata,inline" protobuf:"bytes,1,name=metadata"`
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

// Metadata represents data about a build step
type Metadata struct {
	// Name of the build step
	Name string `json:"name" protobuf:"bytes,1,name=name"`
	// Labels for the step
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,2,opt,name=labels"`
	// Annotations for the step
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,3,opt,name=annotations"`
}

// LoginSpec holds the information to log into a registry.
type LoginSpec struct {
	// Registry refers to a OCI image registry
	Registry string `json:"registry" protobuf:"bytes,1,name=registry"`
	Token    string `json:"token" protobuf:"bytes,1,name=token"`
	// Creds refer to credentials required to log into the registry
	Creds RegistryCreds `json:"creds" protobuf:"bytes,2,name=creds"`
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

// Represents build image metadata
type ImageMetadata struct {
	// BuildFile is the path to the buildfile
	BuildFile string
	// ContextDirectory is the path to the build context
	ContextDirectory string
}
