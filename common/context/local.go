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

package context

import (
	"fmt"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// LocalContext stores the path for your local build context, implements the ContextReader interface
type LocalContext struct {
	// ContextPath is the path to your build context
	ContextPath string `json:"contextPath" protobuf:"bytes,1,opt,name=contextPath"`
}

func (ctx LocalContext) Read() (io.ReadCloser, error) {
	fullPath := fmt.Sprintf("%s/docker-ctx.tar", ctx.ContextPath)

	if ctx.ContextPath == "" {
		return nil, errors.New("cannot have empty contextPath: specify . for current directory")
	}

	defer func(){
		if r := recover(); r != nil {
			logrus.Warnln("panic in context read, recovered for cleanup")
		}
		if err := os.Remove(fullPath); err != nil {
			logrus.WithError(err).Errorln("error cleaning up context file")
		}
	}()

	if err := archiver.Archive([]string{ctx.ContextPath + "/"}, fullPath); err != nil {
		logrus.WithError(err).Errorln("error in building context...")
		return nil, err
	}

	reader, err := os.Open(fullPath)
	if err != nil {
		logrus.WithError(err).Errorln("error in opening docker context...")
		return nil, err
	}

	return reader, nil
}
