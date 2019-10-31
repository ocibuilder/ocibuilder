### ocictl login

- Login to the registry

```
ocictl login -d <PATH_TO_FILE> --builder=docker
ocictl login -d <PATH_TO_FILE> --builder=buildah
```

under the hood, it returns the output of `docker login -u testuser -p testpassword <registry>` or/and `buildah login -u testuser -p testpassword <registry>`. `<registry>` value is fetched from `login.yaml` or `ocibuilder.yaml`

### ocictl build

- Build the image via docker or buildah via ocictl

(Build an image using local Dockerfiles)

```
ocictl build -n <BUILD_NAME> -d <PATH_TO_FILE> --builder=docker
ocictl build -n <BUILD_NAME> -d <PATH_TO_FILE> --builder=buildah
```

under the hood, it returns the output of `docker build -t <image_name> .` and `buildah bud -t <image_name> .` and builds the image. (or could be `docker build -f <path-to-Dockerfile> .` and `buildah bud -f <path-to-Dockerfile> .` and builds the image). `<BUILD_NAME>` is name of the build and `<PATH_TO_FILE>` is path to spec file.

### ocictl pull

3] Pull the image via docker or buildah via ocictl

```
ocictl pull -i <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION> -d <PATH_TO_FILE> --builder=docker
ocictl pull -i <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION> -d <PATH_TO_FILE> --builder=buildah
```

under the hood, it returns the output of `docker pull <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION>`or/and `buildah pull <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION>`. `<REGISTRY-NAME>` value should match with the one in `ocibuilder.yaml` or `login.yaml`

### ocictl push

5] Push images via docker or buildah via ocictl

```
ocictl push -d <PATH_TO_FILE> --builder=docker
ocictl push -d <PATH_TO_FILE> --builder=buildah
```

under the hood, it returns the output of `docker push <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION>` or/and `buildah push <REGISTRY-NAME>/<IMAGE-NAME>:<TAG-VERSION>`. `<REGISTRY-NAME>`, `<IMAGE-NAME>` and `<TAG-VERSION>` is fetched from `ocibuilder.yaml` or `push.yaml`

**Note:** Common functions between ocibuilder/docker and ocibuilder/buildah are under `ocibuilder/common/` directory. We are using go client for executing Docker commands (https://github.com/docker/go-docker). There is no client for buildah. We can use `exec` package in go (exec.Run()) for running buildah commands.

### How to use Overlays

The ocibuilder supports yaml overlays which can be applied at runtime using the `--overlay` command line flag and passing in an overlay yaml file.

The overlaying functionality of ocibuilder is built on top of [ytt](https://github.com/k14s/ytt) which makes use of annotations for matching and overlaying yaml on top of a template. The ocibuilder abstracts away these annotations so you should be able to pass a
standard yaml overlay file and if the fields match up will be applied to your template at runtime. If you want to specify your own ytt annotations, you are able to do so by passing a [standard annotated ytt overlay file](https://get-ytt.io/#example:example-overlay-files)

**Note**: if you want to overlay an array element in yaml make sure to prepend the field with a dash (**-**) so that the ytt annotations can be applied correctly.

Example:

**template**

```yaml
build:
  templates:
    - name: template-1
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
  steps:
    - metadata:
        name: go-build
        labels:
          type: build-1
      stages:
        - metadata:
            name: build-env
          base:
            image: golang
            platform: alpine
          template: template-1
      tag: v0.1.0
```

**overlay**

```yaml
build:
  steps:
    - tag: v0.2.0
```

**overlay applied**

```yaml
build:
  templates:
    - name: template-1
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
  steps:
    - metadata:
        name: go-build
        labels:
          type: build-1
      stages:
        - metadata:
            name: build-env
          base:
            image: golang
            platform: alpine
          template: template-1
      tag: v0.2.0
```

### References

Docker: https://github.com/docker

Buildah Commands: https://github.com/containers/buildah

### Contributing:

Please read the [`CONTRIBUTING.md`](./CONTRIBUTING.md) for contributing guidelines.
