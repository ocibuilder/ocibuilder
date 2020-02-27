FROM golang:1.13.8-alpine
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY dist/oci-init .
ENTRYPOINT ["oci-init"]
