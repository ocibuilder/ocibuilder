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
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	cmdcore "github.com/k14s/ytt/pkg/cmd/core"
	cmdtpl "github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/files"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/request"
	"github.com/pkg/errors"
)

// YttOverlay is the struct for handling overlays using ytt library https://github.com/k14s/ytt
type YttOverlay struct {
	// spec is the spec yaml in a []byte
	Spec []byte
	// overlay is the overlay yaml in a []byte
	Path string
}

// Apply applies the overlay on a YttOverlay struct
func (y YttOverlay) Apply() ([]byte, error) {
	if y.Spec == nil {
		return nil, errors.New("spec file is not defined, overlays is currently only supported for ocibuilder.yaml files")
	}

	overlayFile, err := retrieveOverlayFile(y.Path)

	defer func() {
		if r := recover(); r != nil {
			common.Logger.Warnln("panic recovered to execute final cleanup", r)
		}
		if err := overlayFile.Close(); err != nil {
			common.Logger.WithError(err).Errorln("error closing file")
		}
		if err := os.Remove(common.OverlayPath); err != nil {
			if os.IsNotExist(err) {
				return
			}
			common.Logger.WithError(err).Errorln("error removing file")
		}
	}()

	if err != nil {
		return nil, err
	}

	annotatedOverlay := addYttAnnotations(overlayFile)

	if annotatedOverlay == nil {
		overlay, err := ioutil.ReadFile(y.Path)
		if err != nil {
			return nil, err
		}
		annotatedOverlay = overlay
	}
	filesToProcess := []*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("ocibuilder.yaml", y.Spec)),
		files.MustNewFileFromSource(files.NewBytesSource(y.Path, annotatedOverlay)),
	}

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()
	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		return nil, out.Err
	}
	return out.Files[0].Bytes(), nil
}

func retrieveOverlayFile(path string) (io.ReadCloser, error) {

	// Path is not a valid URI, open local overlay.yaml instead
	if _, err := url.ParseRequestURI(path); err != nil {
		overlayFile, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read overlay file")
		}
		return overlayFile, nil
	}

	if err := request.RequestRemote(path, common.OverlayPath, types.AuthConfig{}); err != nil {
		return nil, err
	}
	overlayFile, err := os.Open(common.OverlayPath)
	if err != nil {
		return nil, err
	}
	return overlayFile, nil
}

// addYttAnnotations adds the expected ytt annotations for ytt overlays
// ytt has a somewhat complex way or applying overlays using annotations. This function
// abstracts away from this such that an overlay is applied w/o the need of annotations.
// If however, the first line of the overlay is #@ load("@ytt:overlay", "overlay"),
// we will default to any specific ytt annotations the user has provided
func addYttAnnotations(overlay io.ReadCloser) []byte {
	yttOverlayIdentifier := "#@ load(\"@ytt:overlay\", \"overlay\")"
	annotatedOverlay := "#@ load(\"@ytt:overlay\", \"overlay\")\n\n#@overlay/match by=overlay.all\n---"

	var tempSegment []string

	addTempToAnnotate := func() {
		annotatedOverlay = annotatedOverlay + "\n" + strings.Join(tempSegment, "\n")
	}
	scanner := bufio.NewScanner(overlay)
	for idx := 0; scanner.Scan(); {
		if idx == 0 && strings.TrimSpace(scanner.Text()) == yttOverlayIdentifier {
			return nil
		}
		if strings.TrimSpace(scanner.Text()) == "- metadata:" {
			addTempToAnnotate()
			tempSegment = nil
		}
		if strings.Contains(scanner.Text(), "overlay:") {
			annotation := retrieveAnnotation(scanner.Text())

			tempSegment = append([]string{annotation}, tempSegment...)
			addTempToAnnotate()
			tempSegment = nil
		}
		tempSegment = append(tempSegment, scanner.Text())
		idx++
	}
	addTempToAnnotate()
	return []byte(annotatedOverlay)
}

func retrieveAnnotation(overlayLine string) string {
	overlayLabel := strings.TrimPrefix(strings.TrimSpace(overlayLine), "overlay:")
	annotationTemplate := "#@overlay/match by=overlay.subset({\"metadata\":{\"labels\":{\"overlay\":\"%s\"}}})"
	return fmt.Sprintf(annotationTemplate, strings.TrimSpace(overlayLabel))
}
