package buildah

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/command"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestClient_ImageBuild(t *testing.T) {
	execute = func(cmd *command.Command) (io.ReadCloser, io.ReadCloser, error) {
		assert.Equal(t, &expectedBuildCommand, cmd)
		return nil, nil, nil
	}
	_, err := cli.ImageBuild(ociBuildOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImagePull(t *testing.T) {
	execute = func(cmd *command.Command) (io.ReadCloser, io.ReadCloser, error) {
		assert.Equal(t, &expectedPullCommand, cmd)
		return nil, nil, nil
	}

	_, err := cli.ImagePull(ociPullOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImagePush(t *testing.T) {
	execute = func(cmd *command.Command) (io.ReadCloser, io.ReadCloser, error) {
		assert.Equal(t, &expectedPushCommand, cmd)
		return nil, nil, nil
	}

	_, err := cli.ImagePush(ociPushOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImageRemove(t *testing.T) {
	execute = func(cmd *command.Command) (io.ReadCloser, io.ReadCloser, error) {
		assert.Equal(t, &expectedRemoveCommand, cmd)
		return nil, nil, nil
	}

	_, err := cli.ImageRemove(ociRemoveOptions)
	assert.Equal(t, nil, err)
}

func TestClient_RegistryLogin(t *testing.T) {
	execute = func(cmd *command.Command) (io.ReadCloser, io.ReadCloser, error) {
		assert.Equal(t, &expectedLoginCommand, cmd)
		return nil, nil, nil
	}

	_, err := cli.RegistryLogin(ociLoginOptions)
	assert.Equal(t, nil, err)
}

func TestClient_GenerateAuthRegistryString(t *testing.T) {
	authString := cli.GenerateAuthRegistryString(authConfig)
	assert.Equal(t, "user:pass", authString)
}

var cli = Client{
	Logger: util.GetLogger(true),
}

var ociBuildOptions = v1alpha1.OCIBuildOptions{
	Ctx:           context.Background(),
	ContextPath:   ".",
	StorageDriver: "vfs",
	ImageBuildOptions: types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		Tags:       []string{"image-name:v0.1.0"},
	},
}

var expectedBuildCommand = command.Builder("buildah").Command("bud").Flags([]command.Flag{
	{Name: "f", Value: "./Dockerfile", Short: true, OmitEmpty: true},
	{Name: "storage-driver", Value: "vfs", Short: false, OmitEmpty: true},
	{Name: "t", Value: "image-name:v0.1.0", Short: true, OmitEmpty: true},
}...).Args(".").Build()

var ociPullOptions = v1alpha1.OCIPullOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePullOptions: types.ImagePullOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var expectedPullCommand = command.Builder("buildah").Command("pull").Flags([]command.Flag{
	{Name: "creds", Value: "this-is-my-auth", Short: false, OmitEmpty: true},
}...).Args("image-name").Build()

var ociPushOptions = v1alpha1.OCIPushOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePushOptions: types.ImagePushOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var expectedPushCommand = command.Builder("buildah").Command("push").Flags([]command.Flag{
	{Name: "creds", Value: "this-is-my-auth", Short: false, OmitEmpty: true},
}...).Args("image-name").Build()

var ociRemoveOptions = v1alpha1.OCIRemoveOptions{
	Image:              "image-name",
	Ctx:                context.Background(),
	ImageRemoveOptions: types.ImageRemoveOptions{},
}

var expectedRemoveCommand = command.Builder("buildah").Command("rmi").Args("image-name").Build()

var ociLoginOptions = v1alpha1.OCILoginOptions{
	Ctx:        context.Background(),
	AuthConfig: authConfig,
}

var expectedLoginCommand = command.Builder("buildah").Command("login").Flags([]command.Flag{
	{Name: "u", Value: "user", Short: true, OmitEmpty: true},
	{Name: "p", Value: "pass", Short: true, OmitEmpty: true},
}...).Args("arts-test-registry").Build()

var authConfig = types.AuthConfig{
	Username:      "user",
	Password:      "pass",
	ServerAddress: "arts-test-registry",
}
