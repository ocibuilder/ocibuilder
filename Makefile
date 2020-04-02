PACKAGE=github.com/beval/beval/provenance
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist
bevalctl_DIR=${DIST_DIR}/bevalctl

VERSION                = $(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE             = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT             = $(shell git rev-parse HEAD)
GENERATE_GROUPS        = vendor/k8s.io/code-generator/generate-groups.sh
DEEPCOPY_GEN           = pkg/apis/beval/v1alpha1/zz_generated.deepcopy.go
DEEPCOPY_SRC           = $(shell find pkg/apis/beval/v1alpha1 | grep -v zz_generated)
CLIENTSET_GEN          = pkg/client/beval/clientset
CLIENTSET_SRC          = $(shell egrep -l -R "^\/\/\s\+genclient" 2>/dev/null | grep -v vendor)
INFORMERS_GEN          = pkg/client/beval/informers
LISTERS_GEN            = pkg/client/beval/listers
HEADER_FILE            = hack/custom-boilerplate.go.txt
beval             = github.com/beval/beval/pkg/client/beval
beval_APIS        = github.com/beval/beval/pkg/apis

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \

define generate
	$(GENERATE_GROUPS) $(1) $(beval) $(beval_APIS) \
	"beval:v1alpha1" --go-header-file $(HEADER_FILE)
endef

.PHONY: beval
beval:
	go build -o ${DIST_DIR}/beval -v .

.PHONY: beval-linux
beval-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make beval

.PHONY: bevalctl
bevalctl: $(bevalctl_DIR)/bevalctl

$(bevalctl_DIR)/bevalctl:
	packr build -v -ldflags '${LDFLAGS}' -o $@ ${CURRENT_DIR}/bevalctl/main.go

.PHONY: bevalctl-linux
bevalctl-linux:
	make clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make bevalctl

.PHONY: bevalctl-mac
bevalctl-mac:
	make clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 make bevalctl

.PHONY: bevalctl-package-linux
bevalctl-package-linux:
	make bevalctl-linux
	cp README.md ./dist/bevalctl
	cp LICENSE ./dist/bevalctl
	cd dist; tar -czvf bevalctl-linux-amd64.tar.gz ./bevalctl

.PHONY: bevalctl-package-mac
bevalctl-package-mac:
	make bevalctl-mac
	cp README.md ./dist/bevalctl
	cp LICENSE ./dist/bevalctl
	cd dist; tar -czvf bevalctl-mac-amd64.tar.gz ./bevalctl

.PHONY: codegen
test:
	go test $(shell go list ./... | grep -v /vendor/ | grep -v /testing/) -race -short -v -coverprofile=coverage.text

.PHONY: lint
lint:
	golangci-lint run

.PHONY: e2e
e2e:
	go test ./testing/e2e -ginkgo.v

.PHONY: clean
clean:
	-rm -rf ${CURRENT_DIR}/dist

.PHONY: clean-gen
clean-gen:
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
