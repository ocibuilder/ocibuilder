package buildah

import (
	"fmt"
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
		{"f", options.Dockerfile, true},
		{"t", options.Tags[0], true},
		{"s", options.StorageDriver, true},
	}

	cmd := common.Builder("buildah").Command("bud").Flags(buildFlags...).Args(options.ContextPath).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing build with command")

	out, err := execute(cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return types.ImageBuildResponse{}, err
	}
	return types.ImageBuildResponse{
		Body: out,
	}, nil
}

func (cli Client) ImagePull(options v1alpha1.OCIPullOptions) (io.ReadCloser, error) {

	pullFlags := []common.Flag{
		// Buildah registry auth in format username[:password]
		{"creds", options.RegistryAuth, false},
	}

	cmd := common.Builder("buildah").Command("pull").Flags(pullFlags...).Args(options.Ref).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing pull with command")

	out, err := execute(cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return nil, err
	}
	return out, nil
}

func (cli Client) ImagePush(options v1alpha1.OCIPushOptions) (io.ReadCloser, error) {

	pushFlags := []common.Flag{
		// Buildah registry auth in format username[:password]
		{"creds", options.RegistryAuth, false},
	}

	cmd := common.Builder("buildah").Command("push").Flags(pushFlags...).Args(options.Ref).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing push with command")

	out, err := execute(cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return nil, err
	}
	return out, nil
}

func (cli Client) ImageRemove(options v1alpha1.OCIRemoveOptions) ([]types.ImageDeleteResponseItem, error) {

	cmd := common.Builder("buildah").Command("rmi").Args(options.Image).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing remove with command")

	_, err := execute(cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return nil, err
	}
	return []types.ImageDeleteResponseItem{
		{
			Deleted: options.Image,
		},
	}, nil
}

func (cli Client) RegistryLogin(options v1alpha1.OCILoginOptions) (registry.AuthenticateOKBody, error) {

	loginFlags := []common.Flag{
		{"u", options.Username, true},
		{"p", options.Password, true},
	}

	cmd := common.Builder("buildah").Command("login").Flags(loginFlags...).Args(options.ServerAddress).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing login with command")

	_, err := execute(cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return registry.AuthenticateOKBody{}, err
	}

	return registry.AuthenticateOKBody{
		Status: "login completed",
	}, nil
}

func (cli Client) GenerateAuthRegistryString(auth types.AuthConfig) string {
	return fmt.Sprintf("%s:%s", auth.Username, auth.Password)
}

// Execute executes the buildah command. This function is mocked in buildah client tests.
var execute = func(cmd common.Command) (io.ReadCloser, error) {
	return cmd.Exec()
}
