/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package overlay

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYttOverlay_Apply(t *testing.T) {
	file, err := os.Open("../testing/dummy/overlay_overlay_test.yaml")
	assert.Equal(t, nil, err)

	yttOverlay := YttOverlay{
		Spec: yamlTplData,
		Overlay: OverlayFile{
			File: file,
			Path: "../testing/dummy/overlay_overlay_test.yaml",
		},
	}
	overlayedSpec, err := yttOverlay.Apply()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedOverlayedSpec, overlayedSpec)
}

func TestYttOverlay_ApplyAnnotated(t *testing.T) {
	file, err := os.Open("../testing/dummy/overlay_overlay_annotated_test.yaml")
	assert.Equal(t, nil, err)

	yttOverlay := YttOverlay{
		Spec: yamlTplData,
		Overlay: OverlayFile{
			File: file,
			Path: "../testing/dummy/overlay_overlay_annotated_test.yaml",
		},
	}
	overlayedSpec, err := yttOverlay.Apply()

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedOverlayedSpec, overlayedSpec)
}

func TestAddYttAnnotations(t *testing.T) {
	file, err := os.Open("../testing/dummy/overlay_overlay_test.yaml")
	assert.Equal(t, nil, err)

	annotatedOverlay := addYttAnnotations(file)
	assert.Equal(t, expectedAnnotatedOverlay, annotatedOverlay)
}

var yamlTplData = []byte(`build:
  templates:
    - name: template-1
      cmd:
        - docker:
            inline:
              - ADD . /src
              - RUN cd /src && go build -o goapp
    - name: template-2
      cmd:
        - docker:
            inline:
              - WORKDIR /app
              - COPY --from=build-env /src/goapp /app/
              - ENTRYPOINT ./goapp
  steps:
    - metadata:
        name: go-build
        labels:
          type: build-1
          overlay: build-1
      stages:
        - metadata:
            name: build-env
            labels:
              stage: stage-1
              overlay: stage-1
          base:
            image: golang
            platform: alpine
          template: template-1
        - metadata:
            name: alpine-stage
            labels:
              stage: stage-2
          base:
            image: alpine
          template: template-2
      tag: v0.1.0
      distroless: false
      cache: false
      purge: false
`)

var expectedAnnotatedOverlay = []byte(`#@ load("@ytt:overlay", "overlay")

#@overlay/match by=overlay.all
---
build:
  steps:
#@overlay/match by=overlay.subset({"metadata":{"labels":{"overlay":"build-1"}}})
    - metadata:
        name: go-service
        labels:
          overlay: build-1
      stages:
#@overlay/match by=overlay.subset({"metadata":{"labels":{"overlay":"stage-1"}}})
        - metadata:
            name: build-env
            labels:
              overlay: stage-1
      tag: v0.2.0`)

var expectedOverlayedSpec = []byte(`build:
  templates:
  - name: template-1
    cmd:
    - docker:
        inline:
        - ADD . /src
        - RUN cd /src && go build -o goapp
  - name: template-2
    cmd:
    - docker:
        inline:
        - WORKDIR /app
        - COPY --from=build-env /src/goapp /app/
        - ENTRYPOINT ./goapp
  steps:
  - metadata:
      name: go-service
      labels:
        type: build-1
        overlay: build-1
    stages:
    - metadata:
        name: build-env
        labels:
          stage: stage-1
          overlay: stage-1
      base:
        image: golang
        platform: alpine
      template: template-1
    - metadata:
        name: alpine-stage
        labels:
          stage: stage-2
      base:
        image: alpine
      template: template-2
    tag: v0.2.0
    distroless: false
    cache: false
    purge: false
`)
