#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE}")/../..
CODEGEN_PKG="$SCRIPT_ROOT/vendor/k8s.io/code-generator"

bash -x ${CODEGEN_PKG}/generate-groups.sh "deepcopy,client,informer,lister" \
				  github.com/ocibuilder/ocibuilder/controller/pkg/client/ocibuilder github.com/ocibuilder/ocibuilder/pkg/apis \
					  "ocibuilder:v1alpha1" \
						  --go-header-file $SCRIPT_ROOT/hack/custom-boilerplate.go.txt

# go run $SCRIPT_ROOT/vendor/k8s.io/gengo/examples/deepcopy-gen/main.go -i github.com/ocibuilder/ocibuilder/pkg/apis --go-header-file $SCRIPT_ROOT/vendor/k8s.io/gengo/boilerplate/boilerplate.go.txt
