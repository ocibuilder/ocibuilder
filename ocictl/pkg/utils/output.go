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

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/ocibuilder/ocibuilder/common"
)

var log = common.GetLogger(false)

// OutputJson streams and formats the output to stdout from returned ReadClosers by docker
// commands.
func OutputJson(ouput io.ReadCloser) error {

	termFd, isTerm := term.GetFdInfo(os.Stdout)

	err := jsonmessage.DisplayJSONMessagesStream(
		ouput,
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
func Output(output io.ReadCloser) error {
	if _, err := io.Copy(os.Stdout, output); err != nil {
		// TODO: this needs to be replaced with a permanent fix, this error is thrown
		// when reaching the end of the reader
		if err.Error() == "read /dev/ptmx: input/output error" {
			log.Infoln("finished reading output")
			return nil
		}
		log.WithError(err).Errorln("error copying output to stdout")
		return err
	}
	return nil
}
