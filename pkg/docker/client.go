package docker

import (
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

type Client struct {
	APIClient client.APIClient
	Logger    *logrus.Logger
}

func (cli Client) ImageBuild(options v1alpha1.OCIBuildOptions) (types.ImageBuildResponse, error) {
	apiCli := cli.APIClient
	return apiCli.ImageBuild(options.Ctx, options.Context, options.ImageBuildOptions)
}

func (cli Client) ImagePull(options v1alpha1.OCIPullOptions) (io.ReadCloser, error) {
	apiCli := cli.APIClient
	return apiCli.ImagePull(options.Ctx, options.Ref, options.ImagePullOptions)
}

func (cli Client) ImagePush(options v1alpha1.OCIPushOptions) (io.ReadCloser, error) {
	apiCli := cli.APIClient
	return apiCli.ImagePush(options.Ctx, options.Ref, options.ImagePushOptions)
}

func (cli Client) ImageRemove(options v1alpha1.OCIRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	apiCli := cli.APIClient
	return apiCli.ImageRemove(options.Ctx, options.Image, options.ImageRemoveOptions)
}

func (cli Client) RegistryLogin(options v1alpha1.OCILoginOptions) (registry.AuthenticateOKBody, error) {
	apiCli := cli.APIClient
	return apiCli.RegistryLogin(options.Ctx, options.AuthConfig)
}

func (cli Client) GenerateAuthRegistryString(auth types.AuthConfig) string {
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error trying to marshall auth config")
	}
	return base64.URLEncoding.EncodeToString(encodedJSON)
}
