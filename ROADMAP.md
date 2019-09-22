# Roadmap

## 0.1.0

* **Buildah and Docker Support** - support for buildah and docker builds
* **Parametrisation and Environment Variables** - support for env variables and parametisation of specifications
* **Multi-stage build** - multi-stage build support for both docker and buildah
* **Purging** - support for purging images after they've been built or after they've been pushed
* **Push** - support for image push to multiple registries
* **Pull** - image pull support for logged in registries
* **Ansible Roles Support** - support for builds using ansible roles

## 0.2.0

* **Image Diff** - support for running container diff on previously built images using [container-diff](https://github.com/GoogleContainerTools/container-diff).
Allows you to easily spot differences in image dependencies.
* **Operator** - kubernetes operator for the ocibuilder allowing you to build images using kubectl.
* **Caching** - caching of image layers
* **Multi Stage Purging** - enhance purging to remove image stages in a multi-stage build
* **Metadata Metrics Storage** - webhook for accessing image build metadata
* **External Configuration** - pull build configuration from external repositories
