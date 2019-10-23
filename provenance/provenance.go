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

// Version information set by link flags during build. We fall back to these sane
// default values when we build outside the Makefile context (e.g. go build or go test).
var (
	// value from VERSION file
	version = "v0.1.0"
	// sha1 from git, output of $(git rev-parse HEAD)
	gitCommit = "$Format:%H$"
	// build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate = "1970-01-01T00:00:00Z"
	// output from `git describe --exact-match --tags HEAD` (if clean tree state)
	gitTag = "v0.1.0"
	// determined from `git status --porcelain`. either 'clean' or 'dirty'
	gitTreeState = ""
)

// Provenance holds information about the build of an executable.
type Provenance struct {
	// Version is the value from VERSION file
	Version string `json:"version"`
	// GitCommit is the output from `git rev-parse HEAD`
	GitCommit string `json:"gitCommit"`
	// BuildDate is output from  `date -u +'%Y-%m-%dT%H:%M:%SZ'`
	BuildDate string `json:"buildDate"`
	// OS holds the operating system name.
	OS string `json:"goOs"`
	// Platform holds architecture name.
	Platform string `json:"goArch"`
	// GitTag refers to tag on a git branch
	GitTag string `json:"gitTag"`
	// GitTreeState is the tree state of git branch/tag
	GitTreeState string `json:"gitTreeState"`
	// Compiler is the go compiler
	Compiler string `json:"compiler"`
	// GoVersion is the version of go language
	GoVersion string
}

// String outputs the version as a string
func (v Provenance) String() string {
	return v.Version
}

// GetProvenance returns an instance of Provenance.
func GetProvenance() Provenance {
	var versionStr string
	if gitCommit != "" && gitTag != "" && gitTreeState == "clean" {
		// if we have a clean tree state and the current commit is tagged,
		// this is an official release.
		versionStr = gitTag
	} else {
		// otherwise formulate a version string based on as much metadata
		// information we have available.
		versionStr = "v" + version
		if len(gitCommit) >= 7 {
			versionStr += "+" + gitCommit[0:7]
			if gitTreeState != "clean" {
				versionStr += ".dirty"
			}
		} else {
			versionStr += "+unknown"
		}
	}
	return Provenance{
		Version:      versionStr,
		BuildDate:    buildDate,
		OS:           runtime.GOOS,
		GitCommit:    gitCommit,
		GitTag:       gitTag,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
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
		v.OS,
		v.Platform); err != nil {
		logrus.WithError(err).Errorln("error printing ocictl verbose version")
	}
}
