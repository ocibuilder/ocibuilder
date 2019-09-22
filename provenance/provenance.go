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

package provenance

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

var (
	version = "unknown"
	// sha1 from git, output of $(git rev-parse HEAD)
	gitCommit = "$Format:%H$"
	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate = "1970-01-01T00:00:00Z"
	goos      = runtime.GOOS
	goarch    = runtime.GOARCH
)

// Provenance holds information about the build of an executable.
type Provenance struct {
	// Version is the value from VERSION file
	Version string `json:"version"`
	// GitCommit is the output from `git rev-parse HEAD`
	GitCommit string `json:"gitCommit"`
	// BuildDate is output from  `date -u +'%Y-%m-%dT%H:%M:%SZ'`
	BuildDate string `json:"buildDate"`
	// GoOs holds OS name.
	GoOs string `json:"goOs"`
	// GoArch holds architecture name.
	GoArch string `json:"goArch"`
}

// GetProvenance returns an instance of Provenance.
func GetProvenance() Provenance {
	return Provenance{
		version,
		gitCommit,
		buildDate,
		goos,
		goarch,
	}
}

func (v Provenance) Print(w io.Writer) {
	if _, err := fmt.Fprintf(w, "%s\n", v.Version); err != nil {
		logrus.WithError(err).Errorln("error printing ocictl version")
	}
}

func (v Provenance) PrintVerbose(w io.Writer) {
	if _, err := fmt.Fprintf(w, "Version: %s, Build Date: %s, Git Commit: %s, OS: %s, arch: %s\n",
		v.Version,
		v.BuildDate,
		v.GitCommit,
		v.GoOs,
		v.GoArch); err != nil {
		logrus.WithError(err).Errorln("error printing ocictl verbose version")
	}
}
