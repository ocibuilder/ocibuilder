package docker

import (
	"context"
	"io"
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestClient_ImageBuild(t *testing.T) {
	_, err := cli.ImageBuild(types.OCIBuildOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImagePull(t *testing.T) {
	_, err := cli.ImagePull(types.OCIPullOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImagePush(t *testing.T) {
	_, err := cli.ImagePush(types.OCIPushOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImageRemove(t *testing.T) {
	_, err := cli.ImageRemove(types.OCIRemoveOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_RegistryLogin(t *testing.T) {
	_, err := cli.RegistryLogin(types.OCILoginOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_GenerateAuthRegistryString(t *testing.T) {
	authString := cli.GenerateAuthRegistryString(authConfig)
	assert.Equal(t, "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6InBhc3MifQ==", authString)
}

var cli = Client{
	Logger:    common.GetLogger(true),
	APIClient: testClient{},
}

var authConfig = dockertypes.AuthConfig{
	Username: "user",
	Password: "pass",
}

func (t testClient) ImageBuild(ctx context.Context, context io.Reader, options dockertypes.ImageBuildOptions) (dockertypes.ImageBuildResponse, error) {
	return dockertypes.ImageBuildResponse{
		Body:   nil,
		OSType: "",
	}, nil
}

func (t testClient) ImagePull(ctx context.Context, ref string, options dockertypes.ImagePullOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (t testClient) ImagePush(ctx context.Context, ref string, options dockertypes.ImagePushOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (t testClient) ImageRemove(ctx context.Context, image string, options dockertypes.ImageRemoveOptions) ([]dockertypes.ImageDeleteResponseItem, error) {
	return nil, nil
}

func (t testClient) RegistryLogin(ctx context.Context, auth dockertypes.AuthConfig) (registry.AuthenticateOKBody, error) {
	return registry.AuthenticateOKBody{
		IdentityToken: "",
		Status:        "",
	}, nil
}

type testClient struct {
	client.APIClient
}
