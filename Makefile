
PACKAGE=github.com/ocibuilder/ocibuilder/provenance
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist

VERSION                = $(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE             = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT             = $(shell git rev-parse HEAD)

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \

# docker image publishing options
DOCKER_PUSH?=true
IMAGE_NAMESPACE?=blackrock
IMAGE_TAG?=v1.0.0

ifeq (${DOCKER_PUSH},true)
ifndef IMAGE_NAMESPACE
$(error IMAGE_NAMESPACE must be set to push images)
endif
endif

ifdef IMAGE_NAMESPACE
IMAGE_PREFIX=${IMAGE_NAMESPACE}/
endif

.PHONY: clean test ocictl

# Proxy
ocibuilder:
	go build -o ${DIST_DIR}/ocibuilder -v .

ocictl:
	packr build -v -ldflags '${LDFLAGS}' -o ${DIST_DIR}/ocictl ${CURRENT_DIR}/ocictl/main.go

ocictl-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr build -v -ldflags '${LDFLAGS}' -o ${DIST_DIR}/ocictl-linux-${VERSION} ${CURRENT_DIR}/ocictl/main.go

ocictl-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 packr build -v -ldflags '${LDFLAGS}' -o ${DIST_DIR}/ocictl-mac-${VERSION} ${CURRENT_DIR}/ocictl/main.go

ocibuilder-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make ocibuilder

ocibuilder-image: ocibuilder-linux
	docker build -t $(IMAGE_PREFIX)db_config_sync:$(IMAGE_TAG) -f Dockerfile .
	@if [ "$(DOCKER_PUSH)" = "true" ] ; then  docker push $(IMAGE_PREFIX)db_config_sync:$(IMAGE_TAG) ; fi

test:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /test/e2e/) -race -short -v

clean:
	-rm -rf ${CURRENT_DIR}/dist
