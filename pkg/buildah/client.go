package buildah

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/command"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Logger *logrus.Logger
}

func (cli Client) ImageBuild(options v1alpha1.OCIBuildOptions) (v1alpha1.OCIBuildResponse, error) {

	buildFlags := []command.Flag{
		{"f", options.Dockerfile, true, true},
		{"storage-driver", options.StorageDriver, false, true},
		{"t", options.Tags[0], true, true},
	}

	cmd := command.Builder("buildah").Command("bud").Flags(buildFlags...).Args(options.ContextPath).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing build with command")

	out, err := execute(&cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return v1alpha1.OCIBuildResponse{}, err
	}
	return v1alpha1.OCIBuildResponse{
		ImageBuildResponse: types.ImageBuildResponse{
			Body: out,
		},
		Exec: &cmd,
	}, nil
}

func (cli Client) ImagePull(options v1alpha1.OCIPullOptions) (v1alpha1.OCIPullResponse, error) {

	pullFlags := []command.Flag{
		// Buildah registry auth in format username[:password]
		{"creds", options.RegistryAuth, false, true},
	}

	cmd := command.Builder("buildah").Command("pull").Flags(pullFlags...).Args(options.Ref).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing pull with command")

	out, err := execute(&cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return v1alpha1.OCIPullResponse{}, err
	}
	return v1alpha1.OCIPullResponse{
		Body: out,
		Exec: &cmd,
	}, nil
}

func (cli Client) ImagePush(options v1alpha1.OCIPushOptions) (v1alpha1.OCIPushResponse, error) {

	pushFlags := []command.Flag{
		// Buildah registry auth in format username[:password]
		{"creds", options.RegistryAuth, false, true},
	}

	cmd := command.Builder("buildah").Command("push").Flags(pushFlags...).Args(options.Ref).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing push with command")

	out, err := execute(&cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return v1alpha1.OCIPushResponse{}, err
	}
	return v1alpha1.OCIPushResponse{
		Body: out,
		Exec: &cmd,
	}, nil
}

func (cli Client) ImageRemove(options v1alpha1.OCIRemoveOptions) (v1alpha1.OCIRemoveResponse, error) {

	cmd := command.Builder("buildah").Command("rmi").Args(options.Image).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing remove with command")

	_, err := execute(&cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return v1alpha1.OCIRemoveResponse{}, err
	}
	return v1alpha1.OCIRemoveResponse{
		Response: []types.ImageDeleteResponseItem{
			{
				Deleted: options.Image,
			},
		},
		Exec: &cmd,
	}, nil
}

func (cli Client) RegistryLogin(options v1alpha1.OCILoginOptions) (v1alpha1.OCILoginResponse, error) {

	loginFlags := []command.Flag{
		{"u", options.Username, true, true},
		{"p", options.Password, true, true},
	}

	cmd := command.Builder("buildah").Command("login").Flags(loginFlags...).Args(options.ServerAddress).Build()
	cli.Logger.WithField("cmd", cmd).Debugln("executing login with command")

	_, err := execute(&cmd)
	if err != nil {
		cli.Logger.WithError(err).Errorln("error building image...")
		return v1alpha1.OCILoginResponse{}, err
	}

	return v1alpha1.OCILoginResponse{
		AuthenticateOKBody: registry.AuthenticateOKBody{
			Status: "login completed",
		},
		Exec: &cmd,
	}, nil
}

func (cli Client) GenerateAuthRegistryString(auth types.AuthConfig) string {
	return fmt.Sprintf("%s:%s", auth.Username, auth.Password)
}

// Execute executes the buildah command. This function is mocked in buildah client tests.
var execute = func(cmd *command.Command) (io.ReadCloser, error) {
	return cmd.Exec()
}
