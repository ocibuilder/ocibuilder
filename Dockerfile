FROM busybox AS ocictl
COPY dist/ocictl /bin/ocictl

FROM ocibuilder/ocibase:0.2.0
COPY --from=binary /bin/ocictl /bin
