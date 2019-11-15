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
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/pkg"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Docker contains configuration required to perform docker operations
type Docker struct {
	// Logger to log stuff
	Logger *logrus.Logger
	// Client is the Docker API client
	Client client.APIClient
}

// Build is used to execute docker build and optionally purge the image after the build
func (d *Docker) Build(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	var buildResponses []io.ReadCloser
	buildOpts, err := pkg.ParseBuildSpec(spec.Build)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse the build specification")
	}
	for _, opt := range buildOpts {
		buildCtx, err := os.Open(opt.BuildContextPath)
		if err != nil {
			return nil, err
		}
		imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)
		dockerOpt := types.ImageBuildOptions{
			Dockerfile: opt.Dockerfile,
			Tags:       []string{imageName},
			Context:    buildCtx,
		}
		buildResponse, err := d.Client.ImageBuild(context.Background(), buildCtx, dockerOpt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to build the image for step %s", opt.Name)
		}
		if err = os.Remove(opt.BuildContextPath + "/" + opt.Dockerfile); err != nil {
			d.Logger.WithError(err).WithField("step", opt.Name).Errorln("error removing generated dockerfile")
		}
		buildResponses = append(buildResponses, buildResponse.Body)
		if opt.Purge {
			res, err := d.Client.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to purge the image for step %s", opt.Name)
			}
			d.Logger.WithField("response", res).Infoln("images purged")
		}
	}
	return buildResponses, nil
}

// Login is used to login to the docker registry/registries
// TODO: review this functionality, login user docker client for go doesn't stick or is inconsistent
func (d *Docker) Login(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	var loginResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		d.Logger.WithField("registry", loginSpec.Registry).Infoln("attempting to login to registry...")
		username, err := pkg.ValidateLoginUsername(loginSpec)
		if err != nil {
			return nil, err
		}
		password, err := pkg.ValidateLoginPassword(loginSpec)
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
		authBody, err := d.Client.RegistryLogin(context.Background(), authCfg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to log into registry %s", loginSpec.Registry)
		}
		loginResponses = append(loginResponses, ioutil.NopCloser(bytes.NewBufferString(authBody.Status)))
	}
	return loginResponses, nil
}

// Pull is used to authenticate with the docker registry and pull an image from the docker registry
func (d *Docker) Pull(spec v1alpha1.OCIBuilderSpec, imageName string) ([]io.ReadCloser, error) {
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
		pullResponse, err := d.Client.ImagePull(context.Background(), registry+imageName, types.ImagePullOptions{
			RegistryAuth: authStr,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "image pull failed for %s", imageName)
		}
		pullResponses = append(pullResponses, pullResponse)
		d.Logger.WithField("image", imageName).Infoln("docker pull has been executed successfully")
	}
	return pullResponses, nil
}

// Push is used to push an image to a docker registry with authentication
func (d Docker) Push(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	var pushResponses []io.ReadCloser
	for _, pushSpec := range spec.Push {
		if err := pkg.ValidatePushSpec(pushSpec); err != nil {
			return nil, err
		}
		pushImageName := fmt.Sprintf("%s/%s:%s", pushSpec.Registry, pushSpec.Image, pushSpec.Tag)
		d.Logger.WithField("name", pushImageName).Infoln("pushing image")
		authString, err := getPushAuthRegistryString(pushSpec.Registry, spec)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get the login credentials for the registry %s", pushSpec.Registry)
		}
		pushResponse, err := d.Client.ImagePush(context.Background(), pushImageName, types.ImagePushOptions{
			RegistryAuth: authString,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to push the image %s", pushImageName)
		}
		if pushSpec.Purge {
			if _, err := d.Client.ImageRemove(context.Background(), pushImageName, types.ImageRemoveOptions{}); err != nil {
				return nil, errors.Wrapf(err, "failed to purge the image %s after push", pushImageName)
			}
			d.Logger.WithField("image", pushImageName).Infoln("image purged")
		}
		pushResponses = append(pushResponses, pushResponse)
	}
	return pushResponses, nil
}

// getPushAuthRegistryString is used to match a push registry with a passed in login specification, returning an auth string
func getPushAuthRegistryString(registry string, spec v1alpha1.OCIBuilderSpec) (string, error) {
	if err := pkg.ValidateLogin(spec); err != nil {
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
	return "", errors.Errorf("no auth credentials matching registry %s found", registry)
}

// encodeAuth is used to generate a base64 encoded auth string to be used in push auth
func encodeAuth(spec v1alpha1.LoginSpec) (string, error) {
	user, err := pkg.ValidateLoginUsername(spec)
	if err != nil {
		return "", err
	}
	pass, err := pkg.ValidateLoginPassword(spec)
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
