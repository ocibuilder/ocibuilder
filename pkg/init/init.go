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

package init

import (
	"io/ioutil"
	"os"

	"github.com/gobuffalo/packr"
	"github.com/sirupsen/logrus"
)

type Initializer struct {
	Box    packr.Box
	Dry    bool
	Logger *logrus.Logger
}

func (i Initializer) Basic() error {
	box := i.Box
	log := i.Logger

	template, err := box.Find("simple_spec_template.yaml")
	if err != nil {
		log.WithError(err).Errorln("error reading in template from docs")
		return err
	}

	if i.Dry {
		if _, err := os.Stdout.Write(template); err != nil {
			log.WithError(err).Errorln("error writing template to stdout")
			return err
		}
	}

	if err := ioutil.WriteFile("ocibuilder.yaml", template, 0777); err != nil {
		log.WithError(err).Errorln("error generating ocibuilder.yaml template file")
		return err
	}

	return nil
}

func (i Initializer) FromDocker(path string) error {
	return nil
}
