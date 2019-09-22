# Configuring Ocibuilder

This document is a reference for ocibuilder v0.1.0 specification keys used in `spec.yaml`.

You can find a complete `spec.yaml` example [here](../../examples/sample_spec.yaml).

---

## Table of Contents
* [`build`](#build)
    * [`templates`](#templates)
        * [`name`](#name)
        * [`cmd`](#cmd)
            * [`ansible`](#ansible)
                * [`galaxy`](#galaxy)
                * [`local`](#local)
            * [`docker`](#docker)
    * [`steps`](#steps)
        * [`imageContext`](#imagecontext)
            * [`localContext`](#localcontext)
        * [`stages`](#stages)
            * [`cmd`](#cmd)
            * [`metadata`](#metadata)
            * [`base`](#base)
        * [`metadata`](#metadata)
* [`login`](#login)
    * [`creds`](#creds)
        * [`plain`](#plain)
        * [`env`](#env)
* [`push`](#push)
* [`params`](#params)

---

### `build`

A build comprises of reusable build templates and a number of build steps

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| templates | *(Array)* v1alpha1.BuildTemplate | Templates are set of build templates that describe steps needed to build a Dockerfile | No |
| steps | *(Array)* v1alpha1.BuildStep | Individual build definitions to run | Yes |


#### `templates`

Templates are reusable build configurations that can be used accross a number of different build steps and are referred to by the name field.
Multiple build commands can be entered which will be executed by ocibuilder.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string | Name of the template | Yes |
| cmd | *(Array)* v1alpha1.BuildTemplateStep | Steps are instructions within a template to build a Dockerfile | Yes |


**Example**
```yaml
  templates:
    - name: go-build-template
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
```

The above example shows a very simple build template, using docker commands which have been passed inline.

#### `name`

Name is the name of the template. This can be referenced across multiple build steps and build stages allowing you to not have to rewrite 
standard build commands across different builds. 

#### `cmd`

Cmd allows you to specify commands that you want to run in your build. This can be standard docker commands, which can be passed 
inline or through a file. Alternatively, you are able to specify an ansible step which you can use to point to a local ansible playbook
or the name and requirements to pull from Ansible galaxy.


| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| ansible | v1alpha1.AnsibleStep | Ansible represents a ansible step within build template steps | No |
| docker | v1alpha1.DockerStep | Docker represents a docker step within build template steps | No |

**Example**
```yaml
  cmd:
    - docker:
        inline:
          - ADD . /src
          - RUN cd /src && go build -o goapp
    - docker:
        path: ./docker-commands.txt
    - ansible:
        galaxy:
          name: my-ansible-role
          requirements: ./requirements.yaml
        local:
          playbook: ./playbook.yaml
          
```

The above example shows all the possible flavours for inputting build commands into ocibuilder.

##### `ansible`

Ansible is used for entering references to local ansible playbooks or ansible roles in Ansible Galaxy. Each ansible step represents an
ansible install as part of the build.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| galaxy | v1alpha1.AnsibleGalaxy | Galaxy contains information to install a ansible role through ansible-galaxy | No |
| local | v1alpha1.AnsibleLocal | Local contains information to install a ansible role through local playbook | No |

##### `galaxy`

AnsibleGalaxy contains ansible role information to be installed at the build

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string | Name of the galaxy role | Yes |
| requirements | string | Requirements refer to the requirements.yaml file | No |

##### `local`

AnsibleLocal contains information to install a ansible role through local playbook

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| playbook | string | Playbook refers to playbook.yaml file | Yes |


##### `docker`

Docker is used for entering docker commands that you want to execute in your build. You have the option of passing commands inline
with each docker command being a string array element.

Alternatively, you can pass in a filepath to a text file which contains your docker commands.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| inline | *(Array)* string | Inline Dockerfile commands | No |
| path | string | Path to a file that contains Dockerfile commands | No |


#### `steps`

Build *steps* are used to configure multiple unique builds with a single `spec.yml`. 

This can be particularly useful when trying to build multiple modules or 
projects in a single repository, allowing you to reuse the `templates` that you have defined in the specification.

Each step is run consecutively by ocibuilder with concurrent build steps running in a future
version of ocibuilder - progress can be tracked [here]().


| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| daemon | boolean | Allows you to specify whether to use the docker daemon as a builder or buildah. Default is true. | No |
| purge | boolean | Purge after built. Defaults is false. | No |
| imageContext | v1alpha1.ImageContext | Specify an image build context for the step | No |
| stages | (*Array*) v1alpha1.Stage | Stages of the build | Yes |
| tag | string | The tag of the built image | No |
| metadata | v1alpha1.Metadata | Build metadata, name, labels, annotations | No |

The purge flag allows you to purge your built images after they've been built. Ensures cleanup and is useful in build pipelines to
prevent the constant persisting of images which will not be used.

**Example**

```yaml
  # Steps are all the individual build steps you want to execute
  steps:
    # Metadata is where you define your final image build name as well as any labels
    - metadata:
        name: my-docker-registry:4555/art/go-service
      stages:
        - metadata:
            name: build-env
          base:
            image: golang
            platform: alpine
          template: go-build-template
      tag: v0.1.0
      purge: false
      daemon: true
      context:
        localContext:
          contextPath: ./go-app
```

##### `imageContext`

ImageContext enables you to specify a build context for each individual build step. For example, if you have
multiple directories and each directory is a separate build step, you are able to specify a particular
directory as the build context that step.

Additionally, an upcoming version of ocibuilder will have support to set external build contexts of S3 buckets and Git paths.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| localContext | context.LocalContext| Local context contains local context information for a build | No |

##### `localContext`

LocalContext holds local context information for an image build

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| contextPath | string | The path to your build context. | Yes |


##### `stages`

Ocibuilder supports [docker multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/) to help 
drastically reduce final image sizes. You are able to define a multi-stage build in each build step.

In each stage you have the option to define build commands under the `cmd` field or pass in a previously specified build template.

A stage also takes in a base image.

Required to pass in either cmd or a build template.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| base | v1alpha1.Base | Refers to parent image for given build stage. | Yes |
| cmd  | (*Array*) v1alpha1.BuildTemplateStep | Cmd refers to a template defined in a stage without a template. | No |
| template | string | Template refers to one of the build templates. | No |
| metadata | v1alpha1.Metadata | Build metadata, name, labels, annotations | No |

**Example**

```yaml
  stages:
    - metadata:
        name: go-binary
      base:
        image: golang
        platform: alpine
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
    - metadata:
        name: alpine-stage
      base:
        image: alpine
      cmd:
        - docker:
            inline:
              - WORKDIR /app
              - COPY --from=go-binary /src/goapp /app/
              - ENTRYPOINT ./goapp
```

In the above example, the first stage of the build uses the golang:alpine base image to build our go binary and is named ``go-binary``.

Our second build stage refers to just our binary built in the first stage with ``--from=go-binary`` and copies this into our
new container image and sets an entrypoint.

More details about multi-stage builds can be found [here](../features/multi-stage-builds.md) 

##### `base`

Base is where you define your base image.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| image | string | The name of the image | Yes |
| platform | string | The specified platform of the image | No |
| tag | string | The tag for the image | No |

**Example**
```yaml
  base:
    image: golang
    platform: alpine
    tag: latest
```

##### `metadata`

This is where any build metadata is defined, including image name, any labels and annotations that you want to specify for your build.

The metadata type is used both in build steps and build stages. 

Within a **build step** name is used to name your final built image, but can be also used to refer to other **build stages** in a
multi-stage build.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| annotations | map | Annotations for your build config | No |
| labels | map | Labels for your build config | No |
| name | string | Name for your build configuration | Yes |

---

### `login`

Login is used to specify credentials necessary to login to an image registry. These credentials are required to push and
pull images from an image registry.

You can specify credentials through plain, with environment variables, through the use of a token or a combination of the above.

[Params](#params) can also be used in conjunction with login credentials to pass in environment variables as a token or plaintext
credentials.

>**NOTE** Using a token may also require you to pass in a username as a credential depending on your image registry.

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| creds | v1alpha1.RegistryCreds | credentials required to log into the registry | Yes |
| registry | string | Registry refers to a OCI image registry | Yes |
| token | string | An image registry token | No |

**Example**
```yaml
login:
  - registry: my-docker-registry:4555
    token: ThisIsMeloGinToken
    creds:
      plain:
        username: art
```

#### `creds`

Holds the specific credentials to login to a given image registry

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| env | v1alpha1.EnvCreds | Env refers to the credentials stored in environment variables | No |
| plain | v1alpha1.PlainCreds | Plain refers to the credentials set inline | No |

##### `env`

Credentials pulled from enviroment variables

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| username | string | Username refers to an env var that holds the username | Yes |
| password | string | Password refers to an env var that holds the password | Yes |

**Example**
```yaml
login:
  env:
    username: REGISTRY_USER
    password: REGISTRY_PASS
```

##### `plain`

Plaintext credentials for image registry login

>**NOTE** It is not recommended to store your credentials in plaintext

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| username | string | Plaintext registry username | Yes |
| password | string | Plaintext registry password | Yes |

**Example**
```yaml
login:
  env:
    username: artsuser
    password: artsp4ss
```

---

### `push`

Push enables you to specify multiple registries to push your image to. Push expects the registry you want to push to
have a corresponding login specification. 

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| image | string | Image to push | Yes |
| registry | string | Registry is the name of the registry | Yes |
| tag | string | Tag version of the image (e.g: v0.1.1) | Yes |

**Example**
```yaml
push:
  - registry: my-docker-registry:4555
    image: art/go-service
    tag: v0.1.0
```

---

### `params`

Params is where you can define parameters, allowing you to replace any value in the specification
with a value or an environment variable.

You specify the destination of the value you want to replace in the `dest` key using dot notation
to access nested elements.

>**NOTE**: A specific array item is referred to by index  in the dest field. For example, if you want to access the first step
element you would have ``steps.0``

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| dest | string | Dest is the destination of the field to replace with the parameter | Yes |
| value | string | Value of the environment variable. | No |
| valueFromEnv | string |  a variable which is to be replaced by an env var | No |

**Example**
```yaml
params:
  # Replaces the value in location build.steps.0.tag with 0.0.3
  - dest: build.steps.0.tag
    value: 0.0.3
  # Replaces the value in location build.steps.0.metadata.name with the environment variable $BUILD_DEV
  - dest: build.steps.0.metadata.name
    valueFromEnv: BUILD_DEV
```