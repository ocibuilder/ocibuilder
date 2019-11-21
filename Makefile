PACKAGE=github.com/ocibuilder/ocibuilder/provenance
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist
OCICTL_DIR=${DIST_DIR}/ocictl

VERSION                = $(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE             = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT             = $(shell git rev-parse HEAD)

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \

.PHONY: clean test ocictl

ocibuilder:
	go build -o ${DIST_DIR}/ocibuilder -v .

ocibuilder-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make ocibuilder

ocictl:
	packr build -v -ldflags '${LDFLAGS}' -o ${OCICTL_DIR}/ocictl ${CURRENT_DIR}/ocictl/main.go

ocictl-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make ocictl

ocictl-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 make ocictl

ocictl-package-build:
	make ocictl-linux
	cp README.md ./dist/ocictl
	cd dist; tar -czvf ocictl-linux-amd64.tar.gz ./ocictl
	make ocictl-mac
	cd dist; tar -czvf ocictl-mac-amd64.tar.gz ./ocictl
	cd dist; rm -rf ./ocictl

test:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /testing/) -race -short -v -coverprofile=coverage.text

lint:
	golangci-lint run

e2e:
	go test testing/e2e

clean:
	-rm -rf ${CURRENT_DIR}/dist

dep:
	dep ensure -v

openapigen:
	hack/update-openapigen.sh

codegen:
	hack/update-codegen.sh
	hack/verify-codegen.sh
