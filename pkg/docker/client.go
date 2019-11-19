package docker

import (
	"encoding/base64"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

type Client struct {
	APIClient client.APIClient
	Logger    *logrus.Logger
}

func (cli Client) ImageBuild(options v1alpha1.OCIBuildOptions) (v1alpha1.OCIBuildResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImageBuild(options.Ctx, options.Context, options.ImageBuildOptions)
	if err != nil {
		return v1alpha1.OCIBuildResponse{}, err
	}
	return v1alpha1.OCIBuildResponse{
		ImageBuildResponse: res,
	}, nil
}

func (cli Client) ImagePull(options v1alpha1.OCIPullOptions) (v1alpha1.OCIPullResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImagePull(options.Ctx, options.Ref, options.ImagePullOptions)
	if err != nil {
		return v1alpha1.OCIPullResponse{}, err
	}
	return v1alpha1.OCIPullResponse{
		Body: res,
	}, nil
}

func (cli Client) ImagePush(options v1alpha1.OCIPushOptions) (v1alpha1.OCIPushResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImagePush(options.Ctx, options.Ref, options.ImagePushOptions)
	if err != nil {
		return v1alpha1.OCIPushResponse{}, err
	}
	return v1alpha1.OCIPushResponse{
		Body: res,
	}, nil
}

func (cli Client) ImageRemove(options v1alpha1.OCIRemoveOptions) (v1alpha1.OCIRemoveResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.ImageRemove(options.Ctx, options.Image, options.ImageRemoveOptions)
	if err != nil {
		return v1alpha1.OCIRemoveResponse{}, err
	}
	return v1alpha1.OCIRemoveResponse{
		Response: res,
	}, nil
}

func (cli Client) RegistryLogin(options v1alpha1.OCILoginOptions) (v1alpha1.OCILoginResponse, error) {
	apiCli := cli.APIClient
	res, err := apiCli.RegistryLogin(options.Ctx, options.AuthConfig)
	if err != nil {
		return v1alpha1.OCILoginResponse{}, err
	}
	return v1alpha1.OCILoginResponse{
		AuthenticateOKBody: res,
	}, nil
}

func (cli Client) GenerateAuthRegistryString(auth types.AuthConfig) string {
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error trying to marshall auth config")
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}