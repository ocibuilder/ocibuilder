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
	file, err := ioutil.ReadFile("../../testing/dummy/spec_docker_gen_test.yaml")
	assert.Equal(t, nil, err)

	spec, err := dockerGen.Generate()
	assert.Equal(t, nil, err)
	assert.Equal(t, string(file), string(spec))
}
