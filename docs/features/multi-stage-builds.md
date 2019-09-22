## Multi-Stage Builds

Multi Stage builds were introduced in Docker 17.05, allowing you to optimize your image builds and share data between
different builds.

The ocibuilder allows you to easily define build stages in your specification in order to run a multi-stage build with 
either Docker or Buildah as the container builder.

#### Defining a multi-stage build

Multi-stage builds can defined at any build step. Each stage has a name defined under it's *metadata* field which can
be referred to at any following build stage.

For example a simple multi-stage go build can look as follows:

```yaml
steps:
- metadata:
    name: artbegolli/go-service
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

The first stage of the build uses the golang:alpine base image to build our go binary and is named ``go-binary``.

Our second build stage refers to just our binary built in the first stage with ``--from=go-binary`` and copies this into our
new container image and sets an entrypoint.

This ultimately results in a significantly smaller image size by creating a minimal image which just contains our build
artifcat.

#### Future Enhancements

- Context labels defined for each stage

#### Links

Docker Multi-Stage Builds https://docs.docker.com/develop/develop-images/multistage-build/