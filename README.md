# ocibuilder - Build [Open Container Initiative (OCI)](https://www.opencontainers.org/) compliant images the declarative way

## What is the ocibuilder?

The **ocibuilder** offers a command line tool called the **ocictl** to build, push and pull [OCI](https://www.opencontainers.org/) compliant images through declarative specifications, allowing
you to pick between [Buildah](https://github.com/containers/buildah) or [Docker](https://docs.docker.com/) as the container build tool. 

## Features

  * Specify docker or buildah as a build tool.
  * Define multiple builds in single build configuration.
  * Ability to templatize build stages.
  * Multi-stage build support
  * Parameterize build configuration at runtime with environment variable support.
  * Supports [distroless](https://github.com/GoogleContainerTools/distroless) to produce lean images.
  * Supports [ansible roles](https://docs.ansible.com/) as build stage.
  * All basic features like registry login, pulling and pushing images from/to multiple registries.

## Architecture

![architecture](https://github.com/ocibuilder/docs/blob/master/assets/ocibuilder.png)

### Documentation

To learn more about ocibuilder and how to get started, [please go to the complete documentation 
here](https://ocibuilder.github.io/docs/)

Our specification file is documented [here](https://ocibuilder.github.io/docs/specification/specification/).

### Roadmap
Take a look at our roadmap and features in developement [here](https://github.com/ocibuilder/ocibuilder/blob/master/ROADMAP.md)

### Contribute

Please read the [`CONTRIBUTING.md`](./CONTRIBUTING.md) for contributing guidelines. 

### License

Apache License Version 2.0, see [`LICENSE`](https://github.com/ocibuilder/ocibuilder/blob/master/LICENSE)

### References

Docker: https://github.com/docker

Buildah Commands: https://github.com/containers/buildah