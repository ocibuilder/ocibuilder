package buildah

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Logger *logrus.Logger
}

func (cli Client) ImageBuild(options v1alpha1.OCIBuildOptions) (types.ImageBuildResponse, error) {

	buildFlags := []common.Flag{
		{"f", options.Dockerfile},
		{"t", options.Tags[0]},
		{"s", options.StorageDriver},
	}

	cmd := common.Builder("buildah").Command("build").Flags(buildFlags...).Args(options.ContextPath).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing build with command")
	out, err := cmd.Exec()
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return types.ImageBuildResponse{}, err
	}
	return types.ImageBuildResponse{
		Body: out,
	}, nil
}

func (cli Client) ImagePull(options v1alpha1.OCIPullOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (cli Client) ImagePush(options v1alpha1.OCIPushOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (cli Client) ImageRemove(options v1alpha1.OCIRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return nil, nil
}

func (cli Client) RegistryLogin(options v1alpha1.OCILoginOptions) (registry.AuthenticateOKBody, error) {
	return registry.AuthenticateOKBody{}, nil
}
