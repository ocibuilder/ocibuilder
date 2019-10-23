/*
Copyright © 2019 BlackRock Inc.

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

package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/sirupsen/logrus"
)

// Docker is a struct which consists of an instance of logger, docker client and context path
type Docker struct {
	Logger *logrus.Logger
	Client client.APIClient
	Metadata []v1alpha1.ImageMeta
}

// Build is used to execute docker build and optionally purge the image after the build
func (d *Docker) Build(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := d.Logger
	cli := d.Client
	buildOpts, err := common.ParseBuildSpec(spec.Build)

	if err != nil {
		log.WithError(err).Errorln("error in parsing build spec...")
		return nil, err
	}

	var buildResponses []io.ReadCloser
	for _, opt := range buildOpts {

		ctx, err := common.ReadContext(opt.Context)
		if err != nil {
			log.WithError(err).Errorln("error reading image build context")
			continue
		}

		imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)

		dockerOpt := types.ImageBuildOptions{
			Dockerfile: opt.Dockerfile,
			Tags:       []string{imageName},
			Context:    ctx,
		}
		buildResponse, err := cli.ImageBuild(context.Background(), ctx, dockerOpt)
		if err != nil {
			log.WithError(err).Errorln("error building image...")
			continue
		}
		buildResponses = append(buildResponses, buildResponse.Body)

		d.Metadata = append(d.Metadata, v1alpha1.ImageMeta{
			BuildFile: opt.Context.LocalContext.ContextPath + "/" + opt.Dockerfile,
		})

		if opt.Purge {
			res, err := cli.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{})
			if err != nil {
				log.WithError(err).Errorln("unable to purge image after build")
				return nil, err
			}
			log.WithFields(logrus.Fields{"response": res}).Infoln("images purged")
		}

	}

	return buildResponses, nil
}

// Login is used to login to the docker registry/registries
// TODO: review this functionality, login user docker client for go doesn't stick or is inconsistent
func (d Docker) Login(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := d.Logger
	cli := d.Client

	var loginResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		log.WithFields(logrus.Fields{"registry": loginSpec.Registry}).Infoln("attempting to login to registry")
		username, err := common.ValidateLoginUsername(loginSpec)
		if err != nil {
			return nil, err
		}

		password, err := common.ValidateLoginPassword(loginSpec)
		if err != nil {
			return nil, err
		}

		if loginSpec.Registry == "" {
			return nil, errors.New("no registry has been specified for login")
		}

		authCfg := types.AuthConfig{
			Username:      username,
			Password:      password,
			ServerAddress: loginSpec.Registry,
		}
		authBody, err := cli.RegistryLogin(context.Background(), authCfg)
		if err != nil {
			log.WithError(err).Errorln("failed to login to registry...")
			return nil, err
		}
		loginResponses = append(loginResponses, ioutil.NopCloser(bytes.NewBufferString(authBody.Status)))
	}
	return loginResponses, nil
}

// Pull is used to authenticate with the docker registry and pull an image from the docker registry
func (d Docker) Pull(spec v1alpha1.OCIBuilderSpec, imageName string) ([]io.ReadCloser, error) {
	log := d.Logger
	cli := d.Client

	var pullResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		registry := loginSpec.Registry
		if registry != "" {
			registry = registry + "/"
		}

		authStr, err := encodeAuth(loginSpec)
		if err != nil {
			return nil, err
		}

		pullResponse, err := cli.ImagePull(context.Background(), registry+imageName, types.ImagePullOptions{
			RegistryAuth: authStr,
		})

		if err != nil {
			log.WithError(err).Errorln("failed to pull image")
			return nil, err
		}
		pullResponses = append(pullResponses, pullResponse)
		log.Infoln("docker pull has been executed successfully")
	}
	return pullResponses, nil
}

// Push is used to push an image to a docker registry with authentication
func (d Docker) Push(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := d.Logger
	cli := d.Client

	var pushResponses []io.ReadCloser
	for _, pushSpec := range spec.Push {
		if err := common.ValidatePushSpec(pushSpec); err != nil {
			return nil, err
		}

		pushImageName := fmt.Sprintf("%s/%s:%s", pushSpec.Registry, pushSpec.Image, pushSpec.Tag)
		log.WithField("name:", pushImageName).Infoln("pushing image")

		authString, err := getPushAuthRegistryString(pushSpec.Registry, spec)
		if err != nil {
			log.WithError(err).Errorln("unable to find login spec")
		}

		pushResponse, err := cli.ImagePush(context.Background(), pushImageName, types.ImagePushOptions{
			RegistryAuth: authString,
		})
		if err != nil {
			log.WithError(err).Errorln("failed to push image")
			return nil, err
		}

		if pushSpec.Purge {
			res, err := cli.ImageRemove(context.Background(), pushImageName, types.ImageRemoveOptions{})
			if err != nil {
				log.WithError(err).Errorln("unable to purge image after push")
				return nil, err
			}
			log.WithFields(logrus.Fields{"response": res}).Infoln("images purged")
		}

		pushResponses = append(pushResponses, pushResponse)
		log.Infoln("docker push has been executed successfully")
	}
	return pushResponses, nil
}

// getPushAuthRegistryString is used to match a push registry with a passed in login specification, returning an auth string
func getPushAuthRegistryString(registry string, spec v1alpha1.OCIBuilderSpec) (string, error) {
	if err := common.ValidateLogin(spec); err != nil {
		return "", err
	}
	for _, spec := range spec.Login {
		if spec.Registry == registry {
			authStr, err := encodeAuth(spec)
			if err != nil {
				return "", err
			}
			return authStr, nil
		}
	}
	return "", errors.New("no auth credentials matching registry: " + registry + " found")
}

// encodeAuth is used to generate a base64 encoded auth string to be used in push auth
func encodeAuth(spec v1alpha1.LoginSpec) (string, error) {
	user, err := common.ValidateLoginUsername(spec)
	if err != nil {
		return "", err
	}

	pass, err := common.ValidateLoginPassword(spec)
	if err != nil {
		return "", err
	}

	authConfig := types.AuthConfig{
		Username: user,
		Password: pass,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr, nil
}

func (d Docker) Clean() {
	log := d.Logger
	for _, m := range d.Metadata {
		if m.BuildFile != "" {
			log.WithField("filepath", m.BuildFile).Debugln("attempting to cleanup dockerfile")
			if err := os.Remove(m.BuildFile); err != nil {
				d.Logger.WithError(err).Errorln("error removing generated Dockerfile")
				continue
			}
		}
	}
}
