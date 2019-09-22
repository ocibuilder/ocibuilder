# ocibuilder
## Version: v1alpha1

### Models


#### v1alpha1.AnsibleGalaxy

AnsibleGalaxy contains the information about the role to install through galaxy

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string | Name of the galaxy role | Yes |
| requirements | string | Requirements refer to the requirements.yaml file | No |

#### v1alpha1.AnsibleLocal

AnsibleLocal contains information to install a ansible role through local playbook

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| playbook | string | Playbook refers to playbook.yaml file | Yes |

#### v1alpha1.AnsibleStep

AnsibleStep represents an ansible install  within a build

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| galaxy | [v1alpha1.AnsibleGalaxy](#v1alpha1.ansiblegalaxy) | Galaxy contains information to install a ansible role through ansible-galaxy | No |
| local | [v1alpha1.AnsibleLocal](#v1alpha1.ansiblelocal) | Local contains information to install a ansible role through local playbook | No |

#### v1alpha1.Base

Base

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| image | string |  | Yes |
| platform | string | Platform is the specified platform of the image | No |
| tag | string | Tag is the tag for the image | No |

#### v1alpha1.BuildSpec

BuildSpec represent the build specifications for images

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| steps | [ [v1alpha1.BuildStep](#v1alpha1.buildstep) ] | Steps within a build | Yes |
| templates | [ [v1alpha1.BuildTemplate](#v1alpha1.buildtemplate) ] | Templates are set of build templates that describe steps needed to build a Dockerfile | Yes |

#### v1alpha1.BuildStep

BuildStep represents a step within the build

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cache | boolean | Cache for build Set to false by default | No |
| context | [v1alpha1.ImageContext](#v1alpha1.imagecontext) | Context used for image build. Default looks at current working directory. | No |
| daemon | boolean | Type of the build framework. Defaults to docker | No |
| distroless | boolean | Distroless if set to true generates a distroless image | No |
| git | string | Git url to fetch the project from. | No |
| purge | boolean | Purge the build defaults to false | No |
| stages | [ [v1alpha1.Stage](#v1alpha1.stage) ] | Stages of the build | Yes |
| tag | string | Tag the tag of the build | No |

#### v1alpha1.BuildTemplate

BuildTemplate represents the build template that can shared across different builds

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cmd | [ [v1alpha1.BuildTemplateStep](#v1alpha1.buildtemplatestep) ] | Steps are instructions within a template to build a Dockerfile | Yes |
| name | string | Name of the template | Yes |

#### v1alpha1.BuildTemplateStep

BuildTemplateStep represents a step within build template

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ansible | [v1alpha1.AnsibleStep](#v1alpha1.ansiblestep) | Ansible represents a ansible step within build template steps | No |
| docker | [v1alpha1.DockerStep](#v1alpha1.dockerstep) | Docker represents a docker step within build template steps | No |

#### v1alpha1.Command

Represents a single line in a Dockerfile

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cmd | string | Cmd lowercased command name (e.g `from`) | Yes |
| flags | [ string ] | Flags such as `--from=...` for `COPY`. | Yes |
| json | boolean | Json bool for whether the value is written in json | Yes |
| original | string | Original is the original source line | Yes |
| startLine | integer | StartLine is the original source line number | Yes |
| subCmd | string | SubCmd for ONBUILD only this holds the sub-command | Yes |
| value | [ string ] | Value is the contents of the command (e.g `ubuntu:xenial`) | Yes |

#### v1alpha1.DockerStep

DockerStep represents a step within a build that contains docker commands

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| inline | [ string ] | Inline Dockerfile commands | No |
| path | string | Path to a file that contains Dockerfile commands | No |

#### v1alpha1.EnvCreds

EnvCreds refers to credentials stored in env vars.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| password | string | Password refers to an en var that holds the password | Yes |
| username | string | Username refers to an env var that holds the username | Yes |

#### v1alpha1.ImageBuildArgs

ImageBuildArgs describes the arguments for running a build command

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ansible | [v1alpha1.AnsibleStep](#v1alpha1.ansiblestep) | Ansible step outlines the ansible steps in the build *optional | No |
| dockerfile | string | Dockerfile is the path to the generated Dockerfile | No |
| name | string | Name is the name of the build | Yes |
| tag | string | Tag is the tag of the build | Yes |

#### v1alpha1.ImageContext

ImageContext stores the chosen build context for your build, this can be Local, S3 or Git

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| gitContext | [context.GitContext](#context.gitcontext) | Git context contains information for a build with a git repo as context | No |
| localContext | [context.LocalContext](#context.localcontext) | Local context contains local context information for a build | No |
| s3Context | [context.S3Context](#context.s3context) | S3 context contains information for a build with an S3 bucket as context | No |

#### v1alpha1.LoginSpec

LoginSpec holds the information to log into a registry.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| creds | [v1alpha1.RegistryCreds](#v1alpha1.registrycreds) | Creds refer to credentials required to log into the registry | Yes |
| registry | string | Registry refers to a OCI image registry | Yes |
| token | string |  | Yes |

#### v1alpha1.Metadata

Metadata

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| annotations | object | Annotations for the step | No |
| labels | object | Labels for the step | No |
| name | string | Name of the build step | Yes |

#### v1alpha1.OCIBuilder

OCIBuilder is the definition of a ocibuilder resource

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| annotations | object | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: http://kubernetes.io/docs/user-guide/annotations | No |
| apiVersion | string | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources | No |
| clusterName | string | The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request. | No |
| deletionGracePeriodSeconds | long | Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only. | No |
| finalizers | [ string ] | Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. | No |
| generateName | string | GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.  If this field is specified and the generated name exists, the server will NOT return a 409 - instead, it will either return 201 Created or 500 with Reason ServerTimeout indicating a unique name could not be found in the time allotted, and the client should retry (optionally after the time indicated in the Retry-After header).  Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#idempotency | No |
| generation | long | A sequence number representing a specific generation of the desired state. Populated by the system. Read-only. | No |
| kind | string | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds | No |
| labels | object | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels | No |
| name | string | Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names | No |
| namespace | string | Namespace defines the space within each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.  Must be a DNS_LABEL. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/namespaces | No |
| resourceVersion | string | An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.  Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#concurrency-control-and-consistency | No |
| selfLink | string | SelfLink is a URL representing this object. Populated by the system. Read-only. | No |
| spec | [v1alpha1.OCIBuilderSpec](#v1alpha1.ocibuilderspec) |  | Yes |
| uid | string | UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.  Populated by the system. Read-only. More info: http://kubernetes.io/docs/user-guide/identifiers#uids | No |

#### v1alpha1.OCIBuilderList

OCIBuilderList is the list of OCIBuilder resources.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| apiVersion | string | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources | No |
| items | [ [v1alpha1.OCIBuilder](#v1alpha1.ocibuilder) ] |  | Yes |
| kind | string | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds | No |

#### v1alpha1.OCIBuilderSpec

OCIBuilderSpec represents OCIBuilder specifications.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| build | [v1alpha1.BuildSpec](#v1alpha1.buildspec) | Build represents the build specifications for images | No |
| login | [ [v1alpha1.LoginSpec](#v1alpha1.loginspec) ] | Logins holds information to log into one or more registries | No |
| params | [ [v1alpha1.Param](#v1alpha1.param) ] | Envs are the list of environment variables available to components. | No |
| push | [ [v1alpha1.PushSpec](#v1alpha1.pushspec) ] | Push contains specification to push images to registries | No |

#### v1alpha1.Param

Param represents parameters

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dest | string | Dest is the destination of the field to replace with the parameter | Yes |
| value | string | Value of the environment variable. | No |
| valueFromEnvVariable | string | ValueFromEnvVar is a variable which is to be replaced by an env var | No |

#### v1alpha1.PlainCreds

Plain refers to the credentials set inline

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| password | string |  | Yes |
| username | string |  | Yes |

#### v1alpha1.PushSpec

PushSpec contains the specification to push images to registries

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| image | string | Image to push | Yes |
| registry | string | Registry is the name of the registry | Yes |
| tag | string | Tag version of the image (e.g: v0.1.1) | Yes |
| token | string | Token required for the OCI complaint registry authentication | Yes |
| user | string | User is the name of kubernetes namespace | Yes |

#### v1alpha1.RegistryCreds

RegistryCreds holds the credentials to login into a registry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| env | [v1alpha1.EnvCreds](#v1alpha1.envcreds) | Env refers to the credentials stored in environment variables | No |
| plain | [v1alpha1.PlainCreds](#v1alpha1.plaincreds) | Plain refers to the credentials set inline | No |

#### v1alpha1.Stage

Stages

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| base | [v1alpha1.Base](#v1alpha1.base) | BaseImage refers to parent image for given build stage. | Yes |
| cmd | [ [v1alpha1.BuildTemplateStep](#v1alpha1.buildtemplatestep) ] | Cmd refers to a template defined in a stage without a template. | Yes |
| template | string | Template refers to one of the build templates. | Yes |

### context.LocalContext

LocalContext holds local context information for an image build

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| contextPath | string | The path to your build context. | Yes |
