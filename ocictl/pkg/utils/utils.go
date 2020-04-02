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

package utils

import (
	"io"
	"os"

	"github.com/beval/beval/pkg/util"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
)

var log = util.GetLogger(false)

// OutputJson streams and formats the output to stdout from returned ReadClosers by docker
// commands.
func OutputJson(output io.ReadCloser) error {

	termFd, isTerm := term.GetFdInfo(os.Stdout)

	err := jsonmessage.DisplayJSONMessagesStream(
		output,
		os.Stdout,
		termFd,
		isTerm,
		nil,
	)

	if err != nil {
		log.WithError(err).Errorln("failed to get JSON stream")
		return err
	}

	return nil

}

// Output outputs a readcloser to stdout in a stream without formatting.
func Output(stdout io.ReadCloser, stderr io.ReadCloser) error {
	//TODO: error with premature read |0: file already closed when finished reading out, investigate further
	if _, err := io.Copy(os.Stdout, stderr); err != nil {
		log.WithError(err).Debugln("error copying output from stderr to stdout, could impact response output")
	}

	if _, err := io.Copy(os.Stdout, stdout); err != nil {
		log.WithError(err).Debugln("error copying output from stdout to stdout, could impact response output")
	}
	return nil
}

// Exists is a simple function to check if a string exists in an array
func Exists(a string, list []string) bool {
	for _, el := range list {
		if el == a {
			return true
		}
	}
	return false
}

// DirExists is a simple function to check if a directory exists
func DirExists(dirName string) bool {
	var exists bool
	if info, err := os.Stat(dirName); err != nil && os.IsNotExist(err) {
		exists = false
	} else {
		exists = info != nil && info.IsDir()
	}
	log.Debugf("Directory [ %s ] exists=%t", dirName, exists)
	return exists
}
