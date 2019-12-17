package generate

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerGenerator_Generate(t *testing.T) {
	dockerGen := DockerGenerator{
		ImageName: "test-image",
		Tag:       "v0.1.0",
		Filepath:  "../../testing/dummy/Dockerfile_Test",
	}
	_, err := ioutil.ReadFile("../../testing/dummy/spec_docker_gen_test.yaml")
	assert.Equal(t, nil, err)

	_, err = dockerGen.Generate()
	assert.Equal(t, nil, err)
}
