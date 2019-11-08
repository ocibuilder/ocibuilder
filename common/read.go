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

package common

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
)

// Read is responsible for reading in the specification files, either
// combined in ocibuilder.yaml or separated in login.yaml, build.yaml and push.yaml.
// The passed in OCIBuilderSpec reference is populated
// If a filepath is not specified the current working directory is used
func (r Reader) Read(spec *v1alpha1.OCIBuilderSpec, overlayPath string, filepaths ...string) error {
	log := r.Logger

	dir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Errorln("failed to get current working directory")
		return err
	}
	filepath := strings.Join(filepaths[:], "/")
	if filepath != "" {
		dir = filepath
	}
	log.WithField("filepath", dir+"/ocibuilder.yaml").Debugln("looking for spec.yaml")
	file, err := ioutil.ReadFile(dir + "/ocibuilder.yaml")
	if err != nil {
		log.Infoln("spec file not found, looking for individual specifications...")
		if err := r.readIndividualSpecs(spec, dir); err != nil {
			log.WithError(err).WithField("directory", dir).Errorln("failed to read individual specs")
			return err
		}
	}

	if err = yaml.Unmarshal(file, spec); err != nil {
		log.WithError(err).WithField("directory", dir).Errorln("failed to unmarshal spec at directory")
		return err
	}

	if err := Validate(spec); err != nil {
		log.WithError(err).WithField("directory", dir).Errorln("failed to validate spec at directory")
		return err
	}

	if overlayPath != "" {
		log.WithField("overlayPath", overlayPath).Debugln("overlay path not empty - looking for overlay file")
		file, err = applyOverlay(file, overlayPath)
		if err != nil {
			log.WithError(err).WithField("path", overlayPath).Errorln("failed to apply overlay to spec at path")
			return err
		}
	}

	if spec.Params != nil {
		if err = r.applyParams(file, spec); err != nil {
			log.WithError(err).Errorln("failed to apply params to spec")
			return err
		}
	}

	return nil
}

// readIndividualSpecs reads the individual specifications if a global
// ocibuilder.yaml is not found
func (r Reader) readIndividualSpecs(spec *v1alpha1.OCIBuilderSpec, path string) error {
	var loginSpec []v1alpha1.LoginSpec
	var buildSpec *v1alpha1.BuildSpec
	var pushSpec []v1alpha1.PushSpec

	r.Logger.Debugln("attempting to read individual specs as spec.yaml as not found")
	if file, err := ioutil.ReadFile(path + "/login.yaml"); err == nil {
		if err := yaml.Unmarshal(file, &loginSpec); err != nil {
			log.WithError(err).Errorln("failed to unmarshal login.yaml")
			return err
		}
		spec.Login = loginSpec
	}
	if file, err := ioutil.ReadFile(path + "/build.yaml"); err == nil {
		if err := yaml.Unmarshal(file, &buildSpec); err != nil {
			log.WithError(err).Errorln("failed to unmarshal build.yaml")
			return err
		}
		spec.Build = buildSpec
	}
	if file, err := ioutil.ReadFile(path + "/push.yaml"); err == nil {
		if err := yaml.Unmarshal(file, &pushSpec); err != nil {
			log.WithError(err).Errorln("failed to unmarshal push.yaml")
			return err
		}
		spec.Push = pushSpec
	}
	return nil
}

// applyOverlay applys a ytt overalay to the specification
func applyOverlay(yamlTemplate []byte, overlayPath string) ([]byte, error) {
	file, err := os.Open(overlayPath)
	if err != nil {
		log.WithError(err).Errorln("unable to read overlay file...")
		return nil, err
	}

	yttOverlay := YttOverlay{
		spec: yamlTemplate,
		overlay: OverlayFile{
			path: overlayPath,
			file: file,
		},
	}

	overlayedSpec, err := yttOverlay.Apply()
	if err != nil {
		log.WithError(err).Errorln("unable to apply overlay to spec...")
		return nil, err
	}

	return overlayedSpec, nil
}

func (r Reader) applyParams(yamlObj []byte, spec *v1alpha1.OCIBuilderSpec) error {
	log := r.Logger
	specJSON, err := yaml.YAMLToJSON(yamlObj)
	if err != nil {
		return err
	}

	log.WithField("number", len(spec.Params)).Debugln("found custom params in spec.yaml")
	for _, param := range spec.Params {
		if param.Value != "" {
			log.WithFields(logrus.Fields{
				"value": param.Value,
				"dest":  param.Dest,
			}).Debugln("setting param value at destination")

			if err := ValidateParams(specJSON, param.Dest); err != nil {
				log.WithError(err).WithField("dest", param.Dest).Errorln("Error validating params, check that your param dest is valid")
				return err
			}

			tmp, err := sjson.SetBytes(specJSON, param.Dest, param.Value)
			if err != nil {
				return err
			}
			specJSON = tmp
		}
		if param.ValueFromEnvVariable != "" {
			log.WithFields(logrus.Fields{
				"value": param.Value,
				"dest":  param.Dest,
			}).Debugln("setting param value at destination")

			val := os.Getenv(param.ValueFromEnvVariable)
			if val == "" {
				log.Warn("env variable ", param.ValueFromEnvVariable, " is empty")
			}

			if err := ValidateParams(specJSON, param.Dest); err != nil {
				log.WithError(err).WithField("dest", param.Dest).Errorln("Error validating params, check that your param dest is valid")
				return err
			}
			tmp, err := sjson.SetBytes(specJSON, param.Dest, val)
			if err != nil {
				return err
			}
			specJSON = tmp
		}
	}

	if err := json.Unmarshal(specJSON, spec); err != nil {
		return err
	}
	return nil
}

// ReadContext reads the user supplied context for the image build
func (r Reader) ReadContext(ctx v1alpha1.ImageContext) (io.ReadCloser, error) {

	if ctx.GitContext != nil {
		return nil, errors.New("git context is not supported in this version of the ocibuilder")
	}

	if ctx.S3Context != nil {
		return nil, errors.New("s3 context is not supported in this version of the ocibuilder")
	}

	if ctx.LocalContext != nil {
		return ctx.LocalContext.Read()
	}

	return nil, errors.New("no context has been defined")
}
