FROM golang:latest AS ocictl
COPY dist/ocictl /bin/ocictl

FROM ocibuilder/ocibase:v0.1.0
COPY --from=ocictl /bin/ocictl /bin
