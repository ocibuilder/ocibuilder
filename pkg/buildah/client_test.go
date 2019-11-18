package buildah

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

func TestClient_ImageBuild(t *testing.T) {
	execute = func(cmd common.Command) (io.ReadCloser, error) {
		assert.Equal(t, expectedBuildCommand, cmd)
		return nil, nil
	}
	_, err := cli.ImageBuild(ociBuildOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImagePull(t *testing.T) {
	execute = func(cmd common.Command) (io.ReadCloser, error) {
		assert.Equal(t, expectedPullCommand, cmd)
		return nil, nil
	}

	_, err := cli.ImagePull(ociPullOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImagePush(t *testing.T) {
	execute = func(cmd common.Command) (io.ReadCloser, error) {
		assert.Equal(t, expectedPushCommand, cmd)
		return nil, nil
	}

	_, err := cli.ImagePush(ociPushOptions)
	assert.Equal(t, nil, err)
}

func TestClient_ImageRemove(t *testing.T) {
	execute = func(cmd common.Command) (io.ReadCloser, error) {
		assert.Equal(t, expectedRemoveCommand, cmd)
		return nil, nil
	}

	_, err := cli.ImageRemove(ociRemoveOptions)
	assert.Equal(t, nil, err)
}

func TestClient_RegistryLogin(t *testing.T) {
	execute = func(cmd common.Command) (io.ReadCloser, error) {
		assert.Equal(t, expectedLoginCommand, cmd)
		return nil, nil
	}

	_, err := cli.RegistryLogin(ociLoginOptions)
	assert.Equal(t, nil, err)
}

func TestClient_GenerateAuthRegistryString(t *testing.T) {
	authString := cli.GenerateAuthRegistryString(authConfig)
	assert.Equal(t, "user:pass", authString)
}

var cli = Client{
	Logger: common.GetLogger(true),
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

var expectedBuildCommand = common.Builder("buildah").Command("bud").Flags([]common.Flag{
	{"f", "./Dockerfile", true},
	{"t", "image-name:v0.1.0", true},
	{"s", "vfs", true},
}...).Args(".").Build()

var ociPullOptions = v1alpha1.OCIPullOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePullOptions: types.ImagePullOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var expectedPullCommand = common.Builder("buildah").Command("pull").Flags([]common.Flag{
	{"creds", "this-is-my-auth", false},
}...).Args("image-name").Build()

var ociPushOptions = v1alpha1.OCIPushOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePushOptions: types.ImagePushOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var expectedPushCommand = common.Builder("buildah").Command("push").Flags([]common.Flag{
	{"creds", "this-is-my-auth", false},
}...).Args("image-name").Build()

var ociRemoveOptions = v1alpha1.OCIRemoveOptions{
	Image:              "image-name",
	Ctx:                context.Background(),
	ImageRemoveOptions: types.ImageRemoveOptions{},
}

var expectedRemoveCommand = common.Builder("buildah").Command("rmi").Args("image-name").Build()

var ociLoginOptions = v1alpha1.OCILoginOptions{
	Ctx:        context.Background(),
	AuthConfig: authConfig,
}

var expectedLoginCommand = common.Builder("buildah").Command("login").Flags([]common.Flag{
	{"u", "user", true},
	{"p", "pass", true},
}...).Args("arts-test-registry").Build()

var authConfig = types.AuthConfig{
	Username:      "user",
	Password:      "pass",
	ServerAddress: "arts-test-registry",
}
