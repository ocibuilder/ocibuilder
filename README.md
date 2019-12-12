# ocibuilder - Build [Open Container Initiative (OCI)](https://www.opencontainers.org/) compliant images the declarative way

[![slack](https://img.shields.io/badge/slack-ocibuilder-brightgreen.svg?logo=slack)](https://urldefense.proofpoint.com/v2/url?u=https-3A__join.slack.com_t_ocibuilder_shared-5Finvite_enQtODYwMTczNzE0OTM1LTZmOWE5MzRlYzI5NzUxYmVmYmIwNjcyN2NlYTZjYTU1ZDAzNjQzZjIyYTQ0NDgwMDdjYzIyZTYyYjYzZWVhZmI&d=DwIFAg&c=zUO0BtkCe66yJvAZ4cAvZg&r=I1Kl9jPjMftKY0cYdNWKyldGY-ke5459qStXhg2CQ9Y&m=CVWkWRAvtaFIja4vZQVm0cJv1-09dRStxuVk4Fr-HBI&s=kkJIZPuFqO6U_nDxBwBTLLczhMnHEt6_f6PZBmcVCs4&e=)
[![Go Report Card](https://goreportcard.com/badge/github.com/ocibuilder/ocibuilder)](https://goreportcard.com/report/github.com/ocibuilder/ocibuilder)
[![CircleCI](https://circleci.com/gh/ocibuilder/ocibuilder.svg?style=shield)](https://circleci.com/gh/ocibuilder/ocibuilder)
[![Docs](https://img.shields.io/badge/docs-ocibuilder-56b5f5)](https://ocibuilder.github.io/docs/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

## What is the ocibuilder?

The **ocibuilder** offers a command line tool called the **ocictl** to build, push and pull [OCI](https://www.opencontainers.org/) compliant images through declarative specifications, allowing you to pick between [Buildah](https://github.com/containers/buildah) or [Docker](https://docs.docker.com/) as the container build tool.

<p align="center">
  <img src="https://github.com/ocibuilder/docs/blob/master/assets/oci-gopher.png?raw=true" alt="Logo"/>
</p>

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

You can use the `install.sh` script to install the latest version of `ocictl`:

```bash
curl https://raw.githubusercontent.com/ocibuilder/ocibuilder/master/install.sh | sh
```

This requires `GOPATH` to be set, with bin added to your `PATH`.

The latest images with Buildah and Docker pre-installed alongside the ocictl is available on our
[Dockerhub repository](https://cloud.docker.com/u/ocibuilder/repository/docker/ocibuilder/ocictl).

Read the full [installation guide](https://ocibuilder.github.io/docs/installation/) available in our docs.

## Getting Started

To learn more about the ocibuilder and how to get started take a look at our [quick start](https://github.com/ocibuilder/docs/blob/master/QUICKSTART.md) guide.

## Documentation

View our complete [documentation](https://ocibuilder.github.io/docs/).

The ocibuilder.yaml specification file with all fields available and examples is documented [here](https://ocibuilder.github.io/docs/specification/specification/).

## Roadmap
Take a look at our roadmap and features in developement [here](https://github.com/ocibuilder/ocibuilder/blob/master/ROADMAP.md)

## Contribute

Please read the [`CONTRIBUTING.md`](./CONTRIBUTING.md) for contributing guidelines.

## License

Apache License Version 2.0, see [`LICENSE`](https://github.com/ocibuilder/ocibuilder/blob/master/LICENSE)

## References

Docker: https://github.com/docker

Buildah Commands: https://github.com/containers/buildah
