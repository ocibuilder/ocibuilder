# ocibuilder - Build [Open Container Initiative (OCI)](https://www.opencontainers.org/) compliant images the declarative way

[![CircleCI](https://circleci.com/gh/ocibuilder/ocibuilder.svg?style=shield)](https://circleci.com/gh/ocibuilder/ocibuilder)
[![Docs](https://img.shields.io/badge/docs-ocibuilder-56b5f5)](https://ocibuilder.github.io/docs/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## What is the ocibuilder?

The **ocibuilder** offers a command line tool called the **ocictl** to build, push and pull [OCI](https://www.opencontainers.org/) compliant images through declarative specifications, allowing you to pick between [Buildah](https://github.com/containers/buildah) or [Docker](https://docs.docker.com/) as the container build tool.

![logo](https://github.com/ocibuilder/docs/blob/master/assets/oci-gopher.png)

## Features

* Specify docker or buildah as a build tool.
* Define multiple builds in single build configuration.
* Ability to templatize build stages.
* Multi-stage build support
* Parameterize build configuration at runtime with environment variable support.
* Supports [distroless](https://github.com/GoogleContainerTools/distroless) to produce lean images.
* Supports [ansible roles](https://docs.ansible.com/) as build stage.
* Supports build contexts like Local Filesystem, Git, S3, Google Cloud Storage, Azure Storage Blob, Aliyun OSS
* All basic features like registry login, pulling and pushing images from/to multiple registries.

## Architecture

![architecture](https://github.com/ocibuilder/docs/blob/master/assets/ocibuilder.png)

## Install

Binary downloads of the `ocictl` are available on the [Releases page](https://github.com/ocibuilder/ocibuilder/releases).

Simply unpack the `ocictl` tar and add it to your path

The latest images with Buildah and Docker pre-installed alongside the ocictl is available on our
[Dockerhub repository](https://cloud.docker.com/u/ocibuilder/repository/docker/ocibuilder/ocictl).

Read the full [installation guide](https://github.com/ocibuilder/docs/blob/master/INSTALL.md) available in our docs.

## Getting Started

To learn more about the ocibuilder and how to get started take a look at our [quick start](https://github.com/ocibuilder/docs/blob/master/QUICKSTART.md) guide.

## Documentation

View our complete [documentation](https://ocibuilder.github.io/docs/).

The ocibuilder.yaml specification file documentation with all fields available with examples can be found [here](https://ocibuilder.github.io/docs/specification/specification/).

## Roadmap
Take a look at our roadmap and features in developement [here](https://github.com/ocibuilder/ocibuilder/blob/master/ROADMAP.md)

## Contribute

Please read the [`CONTRIBUTING.md`](./CONTRIBUTING.md) for contributing guidelines.

## License

Apache License Version 2.0, see [`LICENSE`](https://github.com/ocibuilder/ocibuilder/blob/master/LICENSE)

## References

Docker: https://github.com/docker

Buildah Commands: https://github.com/containers/buildah
