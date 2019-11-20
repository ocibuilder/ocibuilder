package docker

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestClient_ImageBuild(t *testing.T) {
	_, err := cli.ImageBuild(v1alpha1.OCIBuildOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImagePull(t *testing.T) {
	_, err := cli.ImagePull(v1alpha1.OCIPullOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImagePush(t *testing.T) {
	_, err := cli.ImagePush(v1alpha1.OCIPushOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_ImageRemove(t *testing.T) {
	_, err := cli.ImageRemove(v1alpha1.OCIRemoveOptions{})
	assert.Equal(t, nil, err)
}

func TestClient_RegistryLogin(t *testing.T) {
	_, err := cli.RegistryLogin(v1alpha1.OCILoginOptions{})
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

var authConfig = types.AuthConfig{
	Username: "user",
	Password: "pass",
}

func (t testClient) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{
		Body:   nil,
		OSType: "",
	}, nil
}

func (t testClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (t testClient) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	return nil, nil
}

func (t testClient) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return nil, nil
}

func (t testClient) RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error) {
	return registry.AuthenticateOKBody{
		IdentityToken: "",
		Status:        "",
	}, nil
}

type testClient struct {
	client.APIClient
}
