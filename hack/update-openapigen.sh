#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(cd $(dirname "$0")/.. ; pwd)
# CODEGEN_PKG=${PROJECT_ROOT}/vendor/k8s.io/code-generator
OPENAPIGEN_PKG=${PROJECT_ROOT}/vendor/k8s.io/kube-openapi
VERSION="v1alpha1"

# Sensor
# go run ${CODEGEN_PKG}/cmd/openapi-gen/main.go \
go run ${OPENAPIGEN_PKG}/cmd/openapi-gen/openapi-gen.go \
    --go-header-file ${PROJECT_ROOT}/hack/custom-boilerplate.go.txt \
    --input-dirs github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/${VERSION} \
    --output-package /github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/${VERSION} \
    $@
