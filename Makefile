PACKAGE=github.com/ocibuilder/ocibuilder/provenance
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist

VERSION                = $(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE             = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT             = $(shell git rev-parse HEAD)
GENERATE_GROUPS        = vendor/k8s.io/code-generator/generate-groups.sh
DEEPCOPY_GEN           = pkg/apis/ocibuilder/v1alpha1/zz_generated.deepcopy.go
DEEPCOPY_SRC           = $(shell find pkg/apis/ocibuilder/v1alpha1 | grep -v zz_generated)
CLIENTSET_GEN          = pkg/client/ocibuilder/clientset
CLIENTSET_SRC          = $(shell egrep -l -R "^\/\/\s\+genclient" 2>/dev/null | grep -v vendor)
INFORMERS_GEN          = pkg/client/ocibuilder/informers
LISTERS_GEN            = pkg/client/ocibuilder/listers
HEADER_FILE            = hack/custom-boilerplate.go.txt
OCIBUILDER             = github.com/ocibuilder/ocibuilder/pkg/client/ocibuilder
OCIBUILDER_APIS        = github.com/ocibuilder/ocibuilder/pkg/apis

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \

define generate
	$(GENERATE_GROUPS) $(1) $(OCIBUILDER) $(OCIBUILDER_APIS) \
	"ocibuilder:v1alpha1" --go-header-file $(HEADER_FILE)
endef

.PHONY: ocibuilder
ocibuilder:
	go build -o ${DIST_DIR}/ocibuilder -v .

.PHONY: ocibuilder-linux
ocibuilder-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make ocibuilder

.PHONY: ocictl
ocictl: $(DIST_DIR)/ocictl

$(DIST_DIR)/ocictl:
	packr build -v -ldflags '${LDFLAGS}' -o $@ ${CURRENT_DIR}/ocictl/main.go

.PHONY: ocictl-linux
ocictl-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make ocictl

.PHONY: ocictl-mac
ocictl-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 make ocictl

.PHONY: ocictl-package-build
ocictl-package-build:
	make ocictl-linux
	tar -czvf dist/ocictl-linux-${VERSION}.tar.gz ./dist/ocictl
	rm ${CURRENT_DIR}/dist/ocictl
	make ocictl-mac
	tar -czvf dist/ocictl-mac-${VERSION}.tar.gz ./dist/ocictl
	rm ${CURRENT_DIR}/dist/ocictl

.PHONY: codegen
test: codegen
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /testing/) -race -short -v -coverprofile=coverage.text

.PHONY: lint
lint:
	golangci-lint run

.PHONY: e2e
e2e:
	go test testing/e2e

.PHONY: clean
clean:
	-rm -rf ${CURRENT_DIR}/dist
	-rm -rf $(DEEPCOPY_GEN)
	-rm -rf $(CLIENTSET_GEN)
	-rm -rf $(INFORMERS_GEN)
	-rm -rf $(LISTERS_GEN)

.PHONY: dep
dep:
	dep ensure -v

.PHONY: openapigen
openapigen:
	hack/update-openapigen.sh

.PHONY: codegen
codegen: generate-deepcopy generate-client generate-lister generate-informer

.PHONY: verify-codegen
verify-codegen: codegen
	hack/verify-codegen.sh

.PHONY: generate-deepcopy
generate-deepcopy: $(DEEPCOPY_GEN)

$(DEEPCOPY_GEN): $(DEEPCOPY_SRC)
	$(call generate, "deepcopy")

.PHONY: generate-client
generate-client:  $(CLIENTSET_GEN)

$(CLIENTSET_GEN):
	$(call generate, "client")

.PHONY: generate-lister
generate-lister: $(LISTERS_GEN)

$(LISTERS_GEN):
	$(call generate, "lister")

.PHONY: generate-informer
generate-informer: $(INFORMERS_GEN)

$(INFORMERS_GEN):
	$(call generate, "informer")
