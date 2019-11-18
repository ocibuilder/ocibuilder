package buildah

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
)

//TODO: complete test stubs
func TestClient_ImageBuild(t *testing.T) {
}

func TestClient_ImagePull(t *testing.T) {
}

func TestClient_ImagePush(t *testing.T) {
}

func TestClient_ImageRemove(t *testing.T) {
}

func TestClient_RegistryLogin(t *testing.T) {
}

func TestClient_GenerateAuthRegistryString(t *testing.T) {
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

var ociPullOptions = v1alpha1.OCIPullOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePullOptions: types.ImagePullOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var ociPushOptions = v1alpha1.OCIPushOptions{
	Ctx: context.Background(),
	Ref: "image-name",
	ImagePushOptions: types.ImagePushOptions{
		RegistryAuth: "this-is-my-auth",
	},
}

var ociRemoveOptions = v1alpha1.OCIRemoveOptions{
	Image:              "image-name",
	Ctx:                context.Background(),
	ImageRemoveOptions: types.ImageRemoveOptions{},
}

var ociLoginOptions = v1alpha1.OCILoginOptions{
	Ctx:        context.Background(),
	AuthConfig: authConfig,
}

var authConfig = types.AuthConfig{
	Username: "user",
	Password: "pass",
}
