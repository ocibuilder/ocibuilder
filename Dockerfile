FROM busybox AS binary
COPY dist/ocictl /bin/ocictl

FROM ocibuildere2e/ocibuilder-base-go:v0.1.0
COPY --from=binary /bin/ocictl /bin
