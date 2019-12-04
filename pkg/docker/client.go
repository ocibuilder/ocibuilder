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

package docker

import (
	"encoding/base64"
	"encoding/json"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/pkg/types"
	"github.com/sirupsen/logrus"
)

// Client represents docker client
type Client struct {
	APIClient client.APIClient
	Logger    *logrus.Logger
}

func (cli Client) ImageBuild(options types.OCIBuildOptions) (types.OCIBuildResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImageBuild(options.Ctx, options.Context, options.ImageBuildOptions)
	if err != nil {
		return types.OCIBuildResponse{}, err
	}
	return types.OCIBuildResponse{
		ImageBuildResponse: res,
	}, nil
}

func (cli Client) ImagePull(options types.OCIPullOptions) (types.OCIPullResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImagePull(options.Ctx, options.Ref, options.ImagePullOptions)
	if err != nil {
		return types.OCIPullResponse{}, err
	}
	return types.OCIPullResponse{
		Body: res,
	}, nil
}

func (cli Client) ImagePush(options types.OCIPushOptions) (types.OCIPushResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImagePush(options.Ctx, options.Ref, options.ImagePushOptions)
	if err != nil {
		return types.OCIPushResponse{}, err
	}
	return types.OCIPushResponse{
		Body: res,
	}, nil
}

func (cli Client) ImageRemove(options types.OCIRemoveOptions) (types.OCIRemoveResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImageRemove(options.Ctx, options.Image, options.ImageRemoveOptions)
	if err != nil {
		return types.OCIRemoveResponse{}, err
	}
	return types.OCIRemoveResponse{
		Response: res,
	}, nil
}

func (cli Client) RegistryLogin(options types.OCILoginOptions) (types.OCILoginResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.RegistryLogin(options.Ctx, options.AuthConfig)
	if err != nil {
		return types.OCILoginResponse{}, err
	}
	return types.OCILoginResponse{
		AuthenticateOKBody: res,
	}, nil
}

func (cli Client) GenerateAuthRegistryString(auth dockertypes.AuthConfig) string {
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error trying to marshall auth config")
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}
